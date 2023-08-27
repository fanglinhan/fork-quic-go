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
		sender            *bbrSender
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
	// LosePacket := func(number protocol.PacketNumber) {
	// 	sender.OnPacketLost(number, maxDatagramSize, bytesInFlight)
	// 	bytesInFlight -= maxDatagramSize
	// }

	SendAvailableSendWindow := func() int { return SendAvailableSendWindowLen(maxDatagramSize) }
	AckNPackets := func(n int) { AckNPacketsLen(n, maxDatagramSize) }
	LoseNPackets := func(n int) { LoseNPacketsLen(n, maxDatagramSize) }

	BeforeEach(func() {
		bytesInFlight = 0
		packetNumber = 1
		ackedPacketNumber = 0
		clock = mockClock{}
		rttStats = utils.NewRTTStats()
		sender = newBbrSender(
			&clock,
			rttStats,
			protocol.InitialPacketSizeIPv4,
			initialCongestionWindowPackets*maxDatagramSize,
			maxCongestionWindow,
			nil,
		)

	})

	It("has the right values at startup", func() {
		// At startup make sure we are at the default.
		Expect(sender.GetCongestionWindow()).To(Equal(defaultWindowTCP))
		// Make sure we can send.
		Expect(sender.TimeUntilSend(0)).To(BeZero())
		Expect(sender.CanSend(bytesInFlight)).To(BeTrue())
		// And that window is un-affected.
		Expect(sender.GetCongestionWindow()).To(Equal(defaultWindowTCP))

		// Fill the send window with data, then verify that we can't send.
		SendAvailableSendWindow()
		Expect(sender.CanSend(bytesInFlight)).To(BeFalse())
	})

	It("paces", func() {
		rttStats.UpdateRTT(10*time.Millisecond, 0, time.Now())
		clock.Advance(time.Hour)
		// Fill the send window with data, then verify that we can't send.
		SendAvailableSendWindow()
		AckNPackets(1)
		delay := sender.TimeUntilSend(bytesInFlight)
		Expect(delay).ToNot(BeZero())
		Expect(delay).ToNot(Equal(utils.InfDuration))
	})

	// It("application limited slow start", func() {
	// 	// Send exactly 10 packets and ensure the CWND ends at 14 packets.
	// 	const numberOfAcks = 5
	// 	// At startup make sure we can send.
	// 	Expect(sender.CanSend(0)).To(BeTrue())
	// 	Expect(sender.TimeUntilSend(0)).To(BeZero())

	// 	SendAvailableSendWindow()
	// 	for i := 0; i < numberOfAcks; i++ {
	// 		AckNPackets(2)
	// 	}
	// 	bytesToSend := sender.GetCongestionWindow()
	// 	// It's expected 2 acks will arrive when the bytes_in_flight are greater than
	// 	// half the CWND.
	// 	Expect(bytesToSend).To(Equal(defaultWindowTCP + maxDatagramSize*2*2))
	// })

	It("slow start packet loss", func() {
		const numberOfAcks = 10
		for i := 0; i < numberOfAcks; i++ {
			// Send our full send window.
			SendAvailableSendWindow()
			AckNPackets(2)
		}
		SendAvailableSendWindow()
		expectedSendWindow := defaultWindowTCP + (maxDatagramSize * 2 * numberOfAcks)
		Expect(sender.GetCongestionWindow()).To(Equal(expectedSendWindow))

		// Lose a packet to exit slow start.
		LoseNPackets(1)
		packetsInRecoveryWindow := expectedSendWindow / maxDatagramSize

		// We should now have fallen out of slow start with a reduced window.
		expectedSendWindow = protocol.ByteCount(float32(expectedSendWindow) * renoBeta)
		Expect(sender.GetCongestionWindow()).To(Equal(expectedSendWindow))

		// Recovery phase. We need to ack every packet in the recovery window before
		// we exit recovery.
		numberOfPacketsInWindow := expectedSendWindow / maxDatagramSize
		AckNPackets(int(packetsInRecoveryWindow))
		SendAvailableSendWindow()
		Expect(sender.GetCongestionWindow()).To(Equal(expectedSendWindow))

		// We need to ack an entire window before we increase CWND by 1.
		AckNPackets(int(numberOfPacketsInWindow) - 2)
		SendAvailableSendWindow()
		Expect(sender.GetCongestionWindow()).To(Equal(expectedSendWindow))

		// Next ack should increase cwnd by 1.
		AckNPackets(1)
		expectedSendWindow += maxDatagramSize
		Expect(sender.GetCongestionWindow()).To(Equal(expectedSendWindow))
	})
})
