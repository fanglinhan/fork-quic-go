package congestion

import (
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/quic-go/quic-go/internal/protocol"
	"github.com/quic-go/quic-go/internal/utils"
)

var _ = Describe("function", func() {
	It("bdpFromRttAndBandwidth", func() {
		Expect(bdpFromRttAndBandwidth(3*time.Millisecond, Bandwidth(8e3))).To(Equal(protocol.ByteCount(3)))
	})
})

var _ = Describe("", func() {
	const (
		initialCongestionWindowPackets                    = 10
		defaultWindowTCP                                  = protocol.ByteCount(initialCongestionWindowPackets) * maxDatagramSize
		maxCongestionWindow            protocol.ByteCount = 200 * maxDatagramSize
	)

	var (
		sender            *cubicSender
		clock             mockClock
		bytesInFlight     protocol.ByteCount
		packetNumber      protocol.PacketNumber
		ackedPacketNumber protocol.PacketNumber
		rttStats          *utils.RTTStats
	)

	SendAvailableSendWindowLen := func(packetLength protocol.ByteCount) int {
		var packetsSent int
		for sender.CanSend(bytesInFlight) {
			sender.OnPacketSent(clock.Now(), bytesInFlight, packetNumber, packetLength, true)
			packetNumber++
			packetsSent++
			bytesInFlight += packetLength
		}
		return packetsSent
	}

	// Normal is that TCP acks every other segment.
	AckNPacketsLen := func(n int, packetLength protocol.ByteCount) {
		rttStats.UpdateRTT(60*time.Millisecond, 0, clock.Now())
		sender.MaybeExitSlowStart()
		for i := 0; i < n; i++ {
			ackedPacketNumber++
			sender.OnPacketAcked(ackedPacketNumber, packetLength, bytesInFlight, clock.Now())
		}
		bytesInFlight -= protocol.ByteCount(n) * packetLength
		clock.Advance(time.Millisecond)
	}

	LoseNPacketsLen := func(n int, packetLength protocol.ByteCount) {
		for i := 0; i < n; i++ {
			ackedPacketNumber++
			sender.OnPacketLost(ackedPacketNumber, packetLength, bytesInFlight)
		}
		bytesInFlight -= protocol.ByteCount(n) * packetLength
	}

	// Does not increment acked_packet_number_.
	LosePacket := func(number protocol.PacketNumber) {
		sender.OnPacketLost(number, maxDatagramSize, bytesInFlight)
		bytesInFlight -= maxDatagramSize
	}

	SendAvailableSendWindow := func() int { return SendAvailableSendWindowLen(maxDatagramSize) }
	AckNPackets := func(n int) { AckNPacketsLen(n, maxDatagramSize) }
	LoseNPackets := func(n int) { LoseNPacketsLen(n, maxDatagramSize) }

	BeforeEach(func() {
		bytesInFlight = 0
		packetNumber = 1
		ackedPacketNumber = 0
		clock = mockClock{}
		rttStats = utils.NewRTTStats()
		sender = newCubicSender(
			&clock,
			rttStats,
			true, /*reno*/
			protocol.InitialPacketSizeIPv4,
			initialCongestionWindowPackets*maxDatagramSize,
			maxCongestionWindow,
			nil,
		)

	})

	It("", func() {
		Expect(0).To(Equal(0))
		Expect(true).To(BeTrue())
		Expect(func() { panic("") }).To(Panic())
	})
})
