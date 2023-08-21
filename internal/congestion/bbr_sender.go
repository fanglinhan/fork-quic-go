package congestion

import (
	"time"

	"github.com/quic-go/quic-go/internal/protocol"
	"github.com/quic-go/quic-go/internal/utils"
	"github.com/quic-go/quic-go/logging"
)

// BbrSender implements BBR congestion control algorithm.  BBR aims to estimate
// the current available Bottleneck Bandwidth and RTT (hence the name), and
// regulates the pacing rate and the size of the congestion window based on
// those signals.
//
// BBR relies on pacing in order to function properly.  Do not use BBR when
// pacing is disabled.
//

const (
	// Constants based on TCP defaults.
	// The minimum CWND to ensure delayed acks don't reduce bandwidth measurements.
	// Does not inflate the pacing rate.
	defaultMinimumCongestionWindow = 4 * protocol.ByteCount(protocol.InitialPacketSizeIPv4)

	// The gain used for the STARTUP, equal to 2/ln(2).
	defaultHighGain = 2.885
	// The newly derived gain for STARTUP, equal to 4 * ln(2)
	derivedHighGain = 2.773
	// The newly derived CWND gain for STARTUP, 2.
	derivedHighCWNDGain = 2.0
)

// The cycle of gains used during the PROBE_BW stage.
var pacingGain = [...]float64{1.25, 0.75, 1.0, 1.0, 1.0, 1.0, 1.0, 1.0}

const (
	// The length of the gain cycle.
	gainCycleLength = len(pacingGain)
	// The size of the bandwidth filter window, in round-trips.
	bandwidthWindowSize = gainCycleLength + 2

	// The time after which the current min_rtt value expires.
	minRttExpiry = 10 * time.Second
	// The minimum time the connection can spend in PROBE_RTT mode.
	probeRttTime = 200 * time.Millisecond
	// If the bandwidth does not increase by the factor of |kStartupGrowthTarget|
	// within |kRoundTripsWithoutGrowthBeforeExitingStartup| rounds, the connection
	// will exit the STARTUP mode.
	startupGrowthTarget                         = 1.25
	roundTripsWithoutGrowthBeforeExitingStartup = int64(3)
)

type bbrMode int

const (
	// Startup phase of the connection.
	bbrModeStartup = iota
	// After achieving the highest possible bandwidth during the startup, lower
	// the pacing rate in order to drain the queue.
	bbrModeDrain
	// Cruising mode.
	bbrModeProbeBw
	// Temporarily slow down sending in order to empty the buffer and measure
	// the real minimum RTT.
	bbrModeProbeRtt
)

// Indicates how the congestion control limits the amount of bytes in flight.
type bbrRecoveryState int

const (
	// Do not limit.
	bbrRecoveryStateNotInRecovery = iota
	// Allow an extra outstanding byte for each byte acknowledged.
	bbrRecoveryStateConservation
	// Allow two extra outstanding bytes for each byte acknowledged (slow
	// start).
	bbrRecoveryStateGrowth
)

type bbrSender struct {
	rttStats *utils.RTTStats

	mode bbrMode

	// Bandwidth sampler provides BBR with the bandwidth measurements at
	// individual points.
	sampler bandwidthSampler

	// The number of the round trips that have occurred during the connection.
	roundTripCount roundTripCount

	// The packet number of the most recently sent packet.
	lastSendPacket protocol.PacketNumber
	// Acknowledgement of any packet after |current_round_trip_end_| will cause
	// the round trip counter to advance.
	currentRoundTripEnd protocol.PacketNumber

	// Number of congestion events with some losses, in the current round.
	numLossEventsInRound uint64

	// Number of total bytes lost in the current round.
	bytesLostInRound protocol.ByteCount

	// The filter that tracks the maximum bandwidth over the multiple recent
	// round-trips.
	maxBandwidth *utils.WindowedFilter[Bandwidth, roundTripCount]

	// Minimum RTT estimate.  Automatically expires within 10 seconds (and
	// triggers PROBE_RTT mode) if no new value is sampled during that period.
	minRtt time.Duration
	// The time at which the current value of |min_rtt_| was assigned.
	minRttTimestamp time.Time

	// The maximum allowed number of bytes in flight.
	congestionWindow protocol.ByteCount

	// The initial value of the |congestion_window_|.
	initialCongestionWindow protocol.ByteCount

	// The largest value the |congestion_window_| can achieve.
	maxCongestionWindow protocol.ByteCount

	// The smallest value the |congestion_window_| can achieve.
	minCongestionWindow protocol.ByteCount

	// The pacing gain applied during the STARTUP phase.
	highGain float64

	// The CWND gain applied during the STARTUP phase.
	highCwndGain float64

	// The pacing gain applied during the DRAIN phase.
	drainGain float64

	// The current pacing rate of the connection.
	pacingRate Bandwidth

	// The gain currently applied to the pacing rate.
	pacingGain float64
	// The gain currently applied to the congestion window.
	congestionWindowGain float64

	// The gain used for the congestion window during PROBE_BW.  Latched from
	// quic_bbr_cwnd_gain flag.
	congestionWindowGainConstant float64
	// The number of RTTs to stay in STARTUP mode.  Defaults to 3.
	numStartupRtts uint64

	// Number of round-trips in PROBE_BW mode, used for determining the current
	// pacing gain cycle.
	cycleCurrentOffset int
	// The time at which the last pacing gain cycle was started.
	lastCycleStart time.Time

	// Indicates whether the connection has reached the full bandwidth mode.
	isAtFullBandwidth bool
	// Number of rounds during which there was no significant bandwidth increase.
	roundsWithoutBandwidthGain int64
	// The bandwidth compared to which the increase is measured.
	bandwidthAtLastRound Bandwidth

	// Set to true upon exiting quiescence.
	exitingQuiescence bool

	// Time at which PROBE_RTT has to be exited.  Setting it to zero indicates
	// that the time is yet unknown as the number of packets in flight has not
	// reached the required value.
	exitProbeRttAt time.Time
	// Indicates whether a round-trip has passed since PROBE_RTT became active.
	probeRttRoundPassed bool

	// Indicates whether the most recent bandwidth sample was marked as
	// app-limited.
	lastSampleIsAppLimited bool
	// Indicates whether any non app-limited samples have been recorded.
	hasNoAppLimitedSample bool

	// Current state of recovery.
	recoveryState bbrRecoveryState
	// Receiving acknowledgement of a packet after |end_recovery_at_| will cause
	// BBR to exit the recovery mode.  A value above zero indicates at least one
	// loss has been detected, so it must not be set back to zero.
	endRecoveryAt protocol.PacketNumber
	// A window used to limit the number of bytes in flight during loss recovery.
	recoveryWindow protocol.ByteCount
	// If true, consider all samples in recovery app-limited.
	isAppLimitedRecovery bool

	// When true, pace at 1.5x and disable packet conservation in STARTUP.
	slowerStartup bool
	// When true, disables packet conservation in STARTUP.
	rateBasedStartup bool

	// When true, add the most recent ack aggregation measurement during STARTUP.
	enableAckAggregationDuringStartup bool
	// When true, expire the windowed ack aggregation values in STARTUP when
	// bandwidth increases more than 25%.
	expireAckAggregationInStartup bool

	// If true, will not exit low gain mode until bytes_in_flight drops below BDP
	// or it's time for high gain mode.
	drainToTarget bool

	// If true, slow down pacing rate in STARTUP when overshooting is detected.
	detectOvershooting bool
	// Bytes lost while detect_overshooting_ is true.
	bytesLostWhileDetectingOvershooting protocol.ByteCount
	// Slow down pacing rate if
	// bytes_lost_while_detecting_overshooting_ *
	// bytes_lost_multiplier_while_detecting_overshooting_ > IW.
	bytesLostMultiplierWhileDetectingOvershooting uint8
	// When overshooting is detected, do not drop pacing_rate_ below this value /
	// min_rtt.
	cwndToCalculateMinPacingRate protocol.ByteCount

	// Max congestion window when adjusting network parameters.
	maxCongestionWindowWithNetworkParametersAdjusted protocol.ByteCount

	// Params.
	maxDatagramSize protocol.ByteCount
}

var (
	_ SendAlgorithm               = &bbrSender{}
	_ SendAlgorithmWithDebugInfos = &bbrSender{}
)

// NewCubicSender makes a new cubic sender
func NewBbrSender(
	clock Clock,
	rttStats *utils.RTTStats,
	initialMaxDatagramSize protocol.ByteCount,
	tracer logging.ConnectionTracer,
) *bbrSender {
	return newBbrSender(
		clock,
		rttStats,
		initialMaxDatagramSize,
		initialCongestionWindow*initialMaxDatagramSize,
		protocol.MaxCongestionWindowPackets*initialMaxDatagramSize,
		tracer,
	)
}

func newBbrSender(
	clock Clock,
	rttStats *utils.RTTStats,
	initialMaxDatagramSize,
	initialCongestionWindow,
	initialMaxCongestionWindow protocol.ByteCount,
	tracer logging.ConnectionTracer,
) *bbrSender {
	b := &bbrSender{}

	return b
}

// TimeUntilSend implements the SendAlgorithm interface.
func (b *bbrSender) TimeUntilSend(bytesInFlight protocol.ByteCount) time.Time {
	return time.Time{}
}

// HasPacingBudget implements the SendAlgorithm interface.
func (b *bbrSender) HasPacingBudget(now time.Time) bool {
	return false
}

// OnPacketSent implements the SendAlgorithm interface.
func (b *bbrSender) OnPacketSent(sentTime time.Time, bytesInFlight protocol.ByteCount, packetNumber protocol.PacketNumber, bytes protocol.ByteCount, isRetransmittable bool) {

}

// CanSend implements the SendAlgorithm interface.
func (b *bbrSender) CanSend(bytesInFlight protocol.ByteCount) bool {
	return false
}

// MaybeExitSlowStart implements the SendAlgorithm interface.
func (b *bbrSender) MaybeExitSlowStart() {

}

// OnPacketAcked implements the SendAlgorithm interface.
func (b *bbrSender) OnPacketAcked(number protocol.PacketNumber, ackedBytes protocol.ByteCount, priorInFlight protocol.ByteCount, eventTime time.Time) {

}

// OnPacketLost implements the SendAlgorithm interface.
func (b *bbrSender) OnPacketLost(number protocol.PacketNumber, lostBytes protocol.ByteCount, priorInFlight protocol.ByteCount) {

}

// OnRetransmissionTimeout implements the SendAlgorithm interface.
func (b *bbrSender) OnRetransmissionTimeout(packetsRetransmitted bool) {

}

// SetMaxDatagramSize implements the SendAlgorithm interface.
func (b *bbrSender) SetMaxDatagramSize(protocol.ByteCount) {

}

// InSlowStart implements the SendAlgorithmWithDebugInfos interface.
func (b *bbrSender) InSlowStart() bool {
	return false
}

// InRecovery implements the SendAlgorithmWithDebugInfos interface.
func (b *bbrSender) InRecovery() bool {
	return false
}

// GetCongestionWindow implements the SendAlgorithmWithDebugInfos interface.
func (b *bbrSender) GetCongestionWindow() protocol.ByteCount {
	return protocol.MaxByteCount
}

// What's the current estimated bandwidth in bytes per second.
func (b *bbrSender) bandwidthEstimate() Bandwidth {
	return Bandwidth(b.maxBandwidth.GetBest())
}

// Returns the current estimate of the RTT of the connection.  Outside of the
// edge cases, this is minimum RTT.
func (b *bbrSender) getMinRtt() time.Duration {
	if b.minRtt != 0 {
		return b.minRtt
	}
	// min_rtt could be available if the handshake packet gets neutered then
	// gets acknowledged. This could only happen for QUIC crypto where we do not
	// drop keys.
	return b.rttStats.MinRTT()
}

// Computes the target congestion window using the specified gain.
func (b *bbrSender) getTargetCongestionWindow(gain float64) protocol.ByteCount {
	bdp := protocol.ByteCount(b.getMinRtt()) * protocol.ByteCount(b.bandwidthEstimate())
	congestionWindow := protocol.ByteCount(gain * float64(bdp))

	// BDP estimate will be zero if no bandwidth samples are available yet.
	if congestionWindow == 0 {
		congestionWindow = protocol.ByteCount(gain * float64(b.initialCongestionWindow))
	}

	return utils.Max[protocol.ByteCount](congestionWindow, b.minCongestionWindow)
}

// The target congestion window during PROBE_RTT.
func (b *bbrSender) probeRttCongestionWindow() protocol.ByteCount {
	return b.minCongestionWindow
}

// bool MaybeUpdateMinRtt(QuicTime now, QuicTime::Delta sample_min_rtt);

// Enters the STARTUP mode.
// void EnterStartupMode(QuicTime now);

// Enters the PROBE_BW mode.
// void EnterProbeBandwidthMode(QuicTime now);

// Updates the round-trip counter if a round-trip has passed.  Returns true if
// the counter has been advanced.
// bool UpdateRoundTripCounter(QuicPacketNumber last_acked_packet);

// Updates the current gain used in PROBE_BW mode.
// void UpdateGainCyclePhase(QuicTime now, QuicByteCount prior_in_flight, bool has_losses);

// Tracks for how many round-trips the bandwidth has not increased
// significantly.
// void CheckIfFullBandwidthReached(const SendTimeState& last_packet_send_state);

// Transitions from STARTUP to DRAIN and from DRAIN to PROBE_BW if
// appropriate.
// void MaybeExitStartupOrDrain(QuicTime now);

// Decides whether to enter or exit PROBE_RTT.
// void MaybeEnterOrExitProbeRtt(QuicTime now, bool is_round_start, bool min_rtt_expired);

// Determines whether BBR needs to enter, exit or advance state of the
// recovery.
// void UpdateRecoveryState(QuicPacketNumber last_acked_packet, bool has_losses, bool is_round_start);

// Updates the ack aggregation max filter in bytes.
// Returns the most recent addition to the filter, or |newly_acked_bytes| if
// nothing was fed in to the filter.
// QuicByteCount UpdateAckAggregationBytes(QuicTime ack_time, QuicByteCount newly_acked_bytes);

// Determines the appropriate pacing rate for the connection.
func (b *bbrSender) calculatePacingRate(bytesLost protocol.ByteCount) {
	if b.bandwidthEstimate() == 0 {
		return
	}

	targetRate := b.pacingGain * float64(b.bandwidthEstimate())
	if b.isAtFullBandwidth {
		b.pacingRate = Bandwidth(targetRate)
		return
	}

	// Pace at the rate of initial_window / RTT as soon as RTT measurements are
	// available.
	if b.pacingRate == 0 && b.rttStats.MinRTT() != 0 {
		b.pacingRate = BandwidthFromDelta(b.initialCongestionWindow, b.rttStats.MinRTT())
		return
	}

	if b.detectOvershooting {
		b.bytesLostWhileDetectingOvershooting += bytesLost
		// Check for overshooting with network parameters adjusted when pacing rate
		// > target_rate and loss has been detected.
		if b.pacingRate > Bandwidth(targetRate) && b.bytesLostWhileDetectingOvershooting > 0 {
			if b.hasNoAppLimitedSample ||
				b.bytesLostWhileDetectingOvershooting*protocol.ByteCount(b.bytesLostMultiplierWhileDetectingOvershooting) > b.initialCongestionWindow {
				// We are fairly sure overshoot happens if 1) there is at least one
				// non app-limited bw sample or 2) half of IW gets lost. Slow pacing
				// rate.
				b.pacingRate = utils.Max(Bandwidth(targetRate), BandwidthFromDelta(b.cwndToCalculateMinPacingRate, b.rttStats.MinRTT()))
				b.bytesLostWhileDetectingOvershooting = 0
				b.detectOvershooting = false
			}
		}
	}

	// Do not decrease the pacing rate during startup.
	b.pacingRate = utils.Max(b.pacingRate, Bandwidth(targetRate))
}

// Determines the appropriate congestion window for the connection.
func (b *bbrSender) calculateCongestionWindow(bytesAcked, excessAcked protocol.ByteCount) {
	if b.mode == bbrModeProbeRtt {
		return
	}

	targetWindow := b.getTargetCongestionWindow(b.congestionWindowGain)
	if b.isAtFullBandwidth {
		// Add the max recently measured ack aggregation to CWND.
		targetWindow += b.sampler.MaxAckHeight()
	} else if b.enableAckAggregationDuringStartup {
		// Add the most recent excess acked.  Because CWND never decreases in
		// STARTUP, this will automatically create a very localized max filter.
		targetWindow += excessAcked
	}

	// Instead of immediately setting the target CWND as the new one, BBR grows
	// the CWND towards |target_window| by only increasing it |bytes_acked| at a
	// time.
	if b.isAtFullBandwidth {
		b.congestionWindow = utils.Min(targetWindow, b.congestionWindow+bytesAcked)
	} else if b.congestionWindow < targetWindow ||
		b.sampler.TotalBytesAcked() < b.initialCongestionWindow {
		// If the connection is not yet out of startup phase, do not decrease the
		// window.
		b.congestionWindow = b.congestionWindow + bytesAcked
	}

	// Enforce the limits on the congestion window.
	b.congestionWindow = utils.Max(b.congestionWindow, b.minCongestionWindow)
	b.congestionWindow = utils.Min(b.congestionWindow, b.maxCongestionWindow)
}

// Determines the appropriate window that constrains the in-flight during recovery.
func (b *bbrSender) calculateRecoveryWindow(bytesAcked, bytesLost, priorInFlight protocol.ByteCount) {
	if b.recoveryState == bbrRecoveryStateNotInRecovery {
		return
	}

	// Set up the initial recovery window.
	if b.recoveryWindow == 0 {
		b.recoveryWindow = priorInFlight + bytesAcked
		b.recoveryWindow = utils.Max[protocol.ByteCount](b.minCongestionWindow, b.recoveryWindow)
		return
	}

	// Remove losses from the recovery window, while accounting for a potential
	// integer underflow.
	if b.recoveryWindow >= bytesLost {
		b.recoveryWindow = b.recoveryWindow - bytesLost
	} else {
		b.recoveryWindow = b.maxDatagramSize
	}

	// In CONSERVATION mode, just subtracting losses is sufficient.  In GROWTH,
	// release additional |bytes_acked| to achieve a slow-start-like behavior.
	if b.recoveryState == bbrRecoveryStateGrowth {
		b.recoveryWindow += bytesAcked
	}

	// Always allow sending at least |bytes_acked| in response.
	b.recoveryWindow = utils.Max[protocol.ByteCount](b.recoveryWindow, priorInFlight+bytesAcked)
	b.recoveryWindow = utils.Max[protocol.ByteCount](b.minCongestionWindow, b.recoveryWindow)
}

// Called right before exiting STARTUP.
// void OnExitStartup(QuicTime now);

// Return whether we should exit STARTUP due to excessive loss.
// bool ShouldExitStartupDueToLoss(const SendTimeState& last_packet_send_state) const;
