package congestion

import (
	"time"

	"github.com/quic-go/quic-go/internal/protocol"
	"github.com/quic-go/quic-go/internal/utils"
	"github.com/quic-go/quic-go/logging"
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

type bbrRecoveryState int

const (
	// Do not limit.
	bbrStateNotInRecovery = iota
	// Allow an extra outstanding byte for each byte acknowledged.
	bbrStateConservation
	// Allow two extra outstanding bytes for each byte acknowledged (slow
	// start).
	bbrStateGrowth
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

func (b *bbrSender) TimeUntilSend(bytesInFlight protocol.ByteCount) time.Time {
	return time.Time{}
}

func (b *bbrSender) HasPacingBudget(now time.Time) bool {
	return false
}

func (b *bbrSender) OnPacketSent(sentTime time.Time, bytesInFlight protocol.ByteCount, packetNumber protocol.PacketNumber, bytes protocol.ByteCount, isRetransmittable bool) {

}

func (b *bbrSender) CanSend(bytesInFlight protocol.ByteCount) bool {
	return false
}

func (b *bbrSender) MaybeExitSlowStart() {

}

func (b *bbrSender) OnPacketAcked(number protocol.PacketNumber, ackedBytes protocol.ByteCount, priorInFlight protocol.ByteCount, eventTime time.Time) {

}

func (b *bbrSender) OnPacketLost(number protocol.PacketNumber, lostBytes protocol.ByteCount, priorInFlight protocol.ByteCount) {

}

func (b *bbrSender) OnRetransmissionTimeout(packetsRetransmitted bool) {

}

func (b *bbrSender) SetMaxDatagramSize(protocol.ByteCount) {

}

func (b *bbrSender) InSlowStart() bool {
	return false
}

func (b *bbrSender) InRecovery() bool {
	return false
}

func (b *bbrSender) GetCongestionWindow() protocol.ByteCount {
	return protocol.MaxByteCount
}
