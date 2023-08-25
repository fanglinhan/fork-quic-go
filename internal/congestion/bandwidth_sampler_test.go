package congestion

import (
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/quic-go/quic-go/internal/protocol"
	"github.com/quic-go/quic-go/internal/utils"
)

var _ = Describe("function", func() {
	It("BytesFromBandwidthAndTimeDelta", func() {
		Expect(
			bytesFromBandwidthAndTimeDelta(
				Bandwidth(80000),
				100*time.Millisecond,
			)).To(Equal(protocol.ByteCount(1000)))
	})

	It("TimeDeltaFromBytesAndBandwidth", func() {
		Expect(timeDeltaFromBytesAndBandwidth(
			protocol.ByteCount(50000),
			Bandwidth(400),
		)).To(Equal(1000 * time.Second))
	})
})

var _ = Describe("MaxAckHeightTracker", func() {
	var (
		tracker               *maxAckHeightTracker
		now                   time.Time
		rtt                   time.Duration
		bandwidth             Bandwidth
		lastSentPacketNumber  protocol.PacketNumber
		lastAckedPacketNumber protocol.PacketNumber

		getRoundTripCount = func() roundTripCount {
			return roundTripCount(now.Sub(time.Time{}) / rtt)
		}

		// Run a full aggregation episode, which is one or more aggregated acks,
		// followed by a quiet period in which no ack happens.
		// After this function returns, the time is set to the earliest point at which
		// any ack event will cause tracker_.Update() to start a new aggregation.
		aggregationEpisode = func(
			aggregationBandwidth Bandwidth,
			aggregationDuration time.Duration,
			bytesPerAck protocol.ByteCount,
			expectNewAggregationEpoch bool,
		) {
			Expect(aggregationBandwidth >= bandwidth).To(BeTrue())
			startTime := now
			aggregationBytes := bytesFromBandwidthAndTimeDelta(aggregationBandwidth, aggregationDuration)
			numAcks := aggregationBytes / bytesPerAck
			Expect(aggregationBytes).To(Equal(numAcks * bytesPerAck))
			timeBetweenAcks := aggregationDuration / time.Duration(numAcks)
			Expect(aggregationDuration).To(Equal(time.Duration(numAcks) * timeBetweenAcks))

			// The total duration of aggregation time and quiet period.
			totalDuration := timeDeltaFromBytesAndBandwidth(aggregationBytes, bandwidth)
			Expect(aggregationBytes).To(Equal(bytesFromBandwidthAndTimeDelta(bandwidth, totalDuration)))

			var lastExtraAcked protocol.ByteCount
			for bytes := protocol.ByteCount(0); bytes < aggregationBytes; bytes += bytesPerAck {
				extraAcked := tracker.Update(
					bandwidth, true, getRoundTripCount(),
					lastSentPacketNumber, lastAckedPacketNumber, now, bytesPerAck)
				// |extra_acked| should be 0 if either
				// [1] We are at the beginning of a aggregation epoch(bytes==0) and the
				//     the current tracker implementation can identify it, or
				// [2] We are not really aggregating acks.
				if (bytes == 0 && expectNewAggregationEpoch) || (aggregationBandwidth == bandwidth) {
					Expect(extraAcked).To(Equal(protocol.ByteCount(0)))
				} else {
					Expect(lastExtraAcked < extraAcked).To(BeTrue())
				}
				now = now.Add(timeBetweenAcks)
				lastExtraAcked = extraAcked
			}

			// Advance past the quiet period.
			now = startTime.Add(totalDuration)
		}
	)

	BeforeEach(func() {
		tracker = newMaxAckHeightTracker(10)
		tracker.SetAckAggregationBandwidthThreshold(float64(1.8))
		tracker.SetStartNewAggregationEpochAfterFullRound(true)

		now = time.Time{}.Add(1 * time.Millisecond)
		rtt = 60 * time.Millisecond
		bandwidth = Bandwidth(10 * 1000 * 8)
		lastSentPacketNumber = protocol.InvalidPacketNumber
		lastAckedPacketNumber = protocol.InvalidPacketNumber
	})

	It("VeryAggregatedLargeAck", func() {
		aggregationEpisode(bandwidth*20, time.Duration(6*time.Millisecond), 1200, true)
		aggregationEpisode(bandwidth*20, time.Duration(6*time.Millisecond), 1200, true)
		now.Add(-1 * time.Millisecond)

		if tracker.AckAggregationBandwidthThreshold() > float64(1.1) {
			aggregationEpisode(bandwidth*20, time.Duration(6*time.Millisecond), 1200, true)
			Expect(tracker.NumAckAggregationEpochs()).To(Equal(uint64(3)))
		} else {
			aggregationEpisode(bandwidth*20, time.Duration(6*time.Millisecond), 1200, false)
			Expect(tracker.NumAckAggregationEpochs()).To(Equal(uint64(2)))
		}
	})

	It("VeryAggregatedSmallAcks", func() {
		aggregationEpisode(bandwidth*20, time.Duration(6*time.Millisecond), 300, true)
		aggregationEpisode(bandwidth*20, time.Duration(6*time.Millisecond), 300, true)
		now.Add(-1 * time.Millisecond)

		if tracker.AckAggregationBandwidthThreshold() > float64(1.1) {
			aggregationEpisode(bandwidth*20, time.Duration(6*time.Millisecond), 300, true)
			Expect(tracker.NumAckAggregationEpochs()).To(Equal(uint64(3)))
		} else {
			aggregationEpisode(bandwidth*20, time.Duration(6*time.Millisecond), 300, false)
			Expect(tracker.NumAckAggregationEpochs()).To(Equal(uint64(2)))
		}
	})

	It("SomewhatAggregatedLargeAck", func() {
		aggregationEpisode(bandwidth*2, time.Duration(50*time.Millisecond), 1000, true)
		aggregationEpisode(bandwidth*2, time.Duration(50*time.Millisecond), 1000, true)
		now.Add(-1 * time.Millisecond)

		if tracker.AckAggregationBandwidthThreshold() > float64(1.1) {
			aggregationEpisode(bandwidth*2, time.Duration(50*time.Millisecond), 1000, true)
			Expect(tracker.NumAckAggregationEpochs()).To(Equal(uint64(3)))
		} else {
			aggregationEpisode(bandwidth*2, time.Duration(50*time.Millisecond), 1000, false)
			Expect(tracker.NumAckAggregationEpochs()).To(Equal(uint64(2)))
		}
	})

	It("SomewhatAggregatedSmallAcks", func() {
		aggregationEpisode(bandwidth*2, time.Duration(50*time.Millisecond), 100, true)
		aggregationEpisode(bandwidth*2, time.Duration(50*time.Millisecond), 100, true)
		now.Add(-1 * time.Millisecond)

		if tracker.AckAggregationBandwidthThreshold() > float64(1.1) {
			aggregationEpisode(bandwidth*2, time.Duration(50*time.Millisecond), 100, true)
			Expect(tracker.NumAckAggregationEpochs()).To(Equal(uint64(3)))
		} else {
			aggregationEpisode(bandwidth*2, time.Duration(50*time.Millisecond), 100, false)
			Expect(tracker.NumAckAggregationEpochs()).To(Equal(uint64(2)))
		}
	})

	It("NotAggregated", func() {
		aggregationEpisode(bandwidth, time.Duration(100*time.Millisecond), 100, true)
		Expect(uint64(2) < tracker.NumAckAggregationEpochs()).To(BeTrue())
	})

	It("StartNewEpochAfterAFullRound", func() {
		lastSentPacketNumber = protocol.PacketNumber(10)
		aggregationEpisode(bandwidth*2, time.Duration(50*time.Millisecond), 100, true)

		lastAckedPacketNumber = protocol.PacketNumber(11)
		// Update with a tiny bandwidth causes a very low expected bytes acked, which
		// in turn causes the current epoch to continue if the |tracker_| doesn't
		// check the packet numbers.
		tracker.Update(bandwidth/10, true, getRoundTripCount(), lastSentPacketNumber, lastAckedPacketNumber, now, 100)
		Expect(tracker.NumAckAggregationEpochs()).To(Equal(uint64(2)))
	})
})

var _ = Describe("BandwidthSampler", func() {
	var (
		now               time.Time
		sampler           *bandwidthSampler
		regularPacketSize protocol.ByteCount
		// samplerAppLimitedAtStart bool
		bytesInFlight          protocol.ByteCount
		maxBandwidth           Bandwidth // Max observed bandwidth from acks.
		estBandwidthUpperBound Bandwidth
		roundTripCount         roundTripCount // Needed to calculate extra_acked.

		// packetsToBytes = func(packetCount int) protocol.ByteCount {
		// 	return protocol.ByteCount(packetCount) * kRegularPacketSize
		// }

		getPacketSize = func(packetNumber protocol.PacketNumber) protocol.ByteCount {
			return sampler.connectionStateMap.GetEntry(packetNumber).size
		}

		getNumberOfTrackedPackets = func() int {
			return sampler.connectionStateMap.NumberOfPresentEntries()
		}

		sendPacketInner = func(
			packetNumber protocol.PacketNumber,
			bytes protocol.ByteCount,
			hasRetransmittableData bool,
		) {
			sampler.OnPacketSent(now, packetNumber, bytes, bytesInFlight, hasRetransmittableData)
			if hasRetransmittableData {
				bytesInFlight += bytes
			}
		}

		sendPacket = func(packetNumber protocol.PacketNumber) {
			sendPacketInner(packetNumber, regularPacketSize, true)
		}

		makeAckedPacket = func(packetNumber protocol.PacketNumber) protocol.AckedPacketInfo {
			return protocol.AckedPacketInfo{
				PacketNumber: packetNumber,
				BytesAcked:   getPacketSize(packetNumber),
				ReceivedTime: now,
			}
		}

		ackPacketInner = func(packetNumber protocol.PacketNumber) bandwidthSample {
			size := getPacketSize(packetNumber)
			bytesInFlight -= size
			ackedPacket := makeAckedPacket(packetNumber)
			sample := sampler.OnCongestionEvent(now, []protocol.AckedPacketInfo{ackedPacket}, nil,
				maxBandwidth, estBandwidthUpperBound, roundTripCount)
			maxBandwidth = utils.Max(maxBandwidth, sample.sampleMaxBandwidth)
			bwSample := newBandwidthSample()
			bwSample.bandwidth = sample.sampleMaxBandwidth
			bwSample.rtt = sample.sampleRtt
			bwSample.stateAtSend = sample.lastPacketSendState
			Expect(bwSample.stateAtSend.isValid).To(BeTrue())
			return *bwSample
		}

		ackPacket = func(packetNumber protocol.PacketNumber) Bandwidth {
			sample := ackPacketInner(packetNumber)
			return sample.bandwidth
		}

		// makeLostPacket = func(packetNumber protocol.PacketNumber) protocol.LostPacketInfo {
		// 	return protocol.LostPacketInfo{
		// 		PacketNumber: packetNumber,
		// 		BytesLost:    getPacketSize(packetNumber),
		// 	}
		// }

		// losePacket = func(packetNumber protocol.PacketNumber) sendTimeState {
		// 	size := getPacketSize(packetNumber)
		// 	bytesInFlight -= size
		// 	lostPacket := makeLostPacket(packetNumber)
		// 	sample := sampler.OnCongestionEvent(now, nil, []protocol.LostPacketInfo{lostPacket},
		// 		maxBandwidth, estBandwidthUpperBound, roundTripCount)

		// 	Expect(sample.lastPacketSendState.isValid).To(BeTrue())
		// 	Expect(sample.sampleMaxBandwidth).To(Equal(Bandwidth(0)))
		// 	Expect(sample.sampleRtt).To(Equal(infRTT))
		// 	return sample.lastPacketSendState
		// }

		// onCongestionEvent = func(ackedPacketNumbers, lostPacketNumbers []protocol.PacketNumber) congestionEventSample {
		// 	ackedPackets := []protocol.AckedPacketInfo{}
		// 	for _, packetNumber := range ackedPacketNumbers {
		// 		ackedPacket := makeAckedPacket(packetNumber)
		// 		ackedPackets = append(ackedPackets, makeAckedPacket(packetNumber))
		// 		bytesInFlight -= ackedPacket.BytesAcked
		// 	}
		// 	lostPackets := []protocol.LostPacketInfo{}
		// 	for _, packetNumber := range lostPacketNumbers {
		// 		lostPacket := makeLostPacket(packetNumber)
		// 		lostPackets = append(lostPackets, lostPacket)
		// 		bytesInFlight -= lostPacket.BytesLost
		// 	}

		// 	sample := sampler.OnCongestionEvent(now, ackedPackets, lostPackets,
		// 		maxBandwidth, estBandwidthUpperBound, roundTripCount)
		// 	maxBandwidth = utils.Max(maxBandwidth, sample.sampleMaxBandwidth)
		// 	return sample
		// }

		// Sends one packet and acks it.  Then, send 20 packets.  Finally, send
		// another 20 packets while acknowledging previous 20.
		// send40PacketsAndAckFirst20 = func(timeBetweenPackets time.Duration) {
		// 	// Send 20 packets at a constant inter-packet time.
		// 	for i := 1; i <= 20; i++ {
		// 		sendPacket(protocol.PacketNumber(i))
		// 		now = now.Add(timeBetweenPackets)
		// 	}

		// 	// Ack packets 1 to 20, while sending new packets at the same rate as
		// 	// before.
		// 	for i := 1; i <= 20; i++ {
		// 		ackPacket(protocol.PacketNumber(i))
		// 		sendPacket(protocol.PacketNumber(i + 20))
		// 		now = now.Add(timeBetweenPackets)
		// 	}
		// }

		testParameters = []struct {
			overestimateAvoidance bool
		}{
			{
				overestimateAvoidance: false,
			},
			{
				overestimateAvoidance: true,
			},
		}

		initial = func(param struct {
			overestimateAvoidance bool
		}) {
			// Ensure that the clock does not start at zero.
			now = time.Time{}.Add(1 * time.Second)
			sampler = newBandwidthSampler(0)
			regularPacketSize = protocol.ByteCount(1280)
			// samplerAppLimitedAtStart = false
			bytesInFlight = protocol.ByteCount(0)
			maxBandwidth = 0
			estBandwidthUpperBound = infBandwidth
			roundTripCount = 0

			if param.overestimateAvoidance {
				sampler.EnableOverestimateAvoidance()
			}
		}
	)

	// Test the sampler in a simple stop-and-wait sender setting.
	It("SendAndWait", func() {
		for _, param := range testParameters {
			initial(param)

			timeBetweenPackets := 10 * time.Millisecond
			expectedBandwidth := Bandwidth(regularPacketSize) * 100 * BytesPerSecond

			// Send packets at the constant bandwidth.
			for i := 1; i < 20; i++ {
				sendPacket(protocol.PacketNumber(i))
				now = now.Add(timeBetweenPackets)
				currentSample := ackPacket(protocol.PacketNumber(i))
				Expect(expectedBandwidth).To(Equal(currentSample))
			}

			// Send packets at the exponentially decreasing bandwidth.
			for i := 20; i < 25; i++ {
				timeBetweenPackets *= 2
				expectedBandwidth /= 2

				sendPacket(protocol.PacketNumber(i))
				now = now.Add(timeBetweenPackets)
				currentSample := ackPacket(protocol.PacketNumber(i))
				Expect(expectedBandwidth).To(Equal(currentSample))
			}
			sampler.RemoveObsoletePackets(protocol.PacketNumber(25))

			Expect(getNumberOfTrackedPackets()).To(Equal(int(0)))
			Expect(bytesInFlight).To(Equal(protocol.ByteCount(0)))
		}
	})

	It("SendTimeState", func() {
		for _, param := range testParameters {
			initial(param)

		}
	})

	It("SendPaced", func() {
		for _, param := range testParameters {
			initial(param)

		}
	})

	It("SendWithLosses", func() {
		for _, param := range testParameters {
			initial(param)

		}
	})

	It("NotCongestionControlled", func() {
		for _, param := range testParameters {
			initial(param)

		}
	})

	It("CompressedAck", func() {
		for _, param := range testParameters {
			initial(param)

		}
	})

	It("ReorderedAck", func() {
		for _, param := range testParameters {
			initial(param)

		}
	})

	It("AppLimited", func() {
		for _, param := range testParameters {
			initial(param)

		}
	})

	It("FirstRoundTrip", func() {
		for _, param := range testParameters {
			initial(param)

		}
	})

	It("RemoveObsoletePackets", func() {
		for _, param := range testParameters {
			initial(param)

		}
	})

	It("NeuterPacket", func() {
		for _, param := range testParameters {
			initial(param)

		}
	})

	It("CongestionEventSampleDefaultValues", func() {
		for _, param := range testParameters {
			initial(param)

		}
	})

	It("TwoAckedPacketsPerEvent", func() {
		for _, param := range testParameters {
			initial(param)

		}
	})

	It("LoseEveryOtherPacket", func() {
		for _, param := range testParameters {
			initial(param)

		}
	})

	It("AckHeightRespectBandwidthEstimateUpperBound", func() {
		for _, param := range testParameters {
			initial(param)

		}
	})
})
