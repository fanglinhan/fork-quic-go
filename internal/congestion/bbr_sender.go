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
	return &bbrSender{}
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
