package congestion

import (
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/quic-go/quic-go/internal/protocol"
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
				now.Add(timeBetweenAcks)
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
			Expect(tracker.numAckAggregationEpochs).To(Equal(uint64(3)))
		} else {
			aggregationEpisode(bandwidth*20, time.Duration(6*time.Millisecond), 1200, false)
			Expect(tracker.numAckAggregationEpochs).To(Equal(uint64(2)))
		}
	})

	It("VeryAggregatedSmallAcks", func() {
		aggregationEpisode(bandwidth*20, time.Duration(6*time.Millisecond), 300, true)
		aggregationEpisode(bandwidth*20, time.Duration(6*time.Millisecond), 300, true)
		now.Add(-1 * time.Millisecond)

		if tracker.AckAggregationBandwidthThreshold() > float64(1.1) {
			aggregationEpisode(bandwidth*20, time.Duration(6*time.Millisecond), 300, true)
			Expect(tracker.numAckAggregationEpochs).To(Equal(uint64(3)))
		} else {
			aggregationEpisode(bandwidth*20, time.Duration(6*time.Millisecond), 300, false)
			Expect(tracker.numAckAggregationEpochs).To(Equal(uint64(2)))
		}
	})

	It("SomewhatAggregatedLargeAck", func() {
		aggregationEpisode(bandwidth*2, time.Duration(50*time.Millisecond), 1000, true)
		aggregationEpisode(bandwidth*2, time.Duration(50*time.Millisecond), 1000, true)
		now.Add(-1 * time.Millisecond)

		if tracker.AckAggregationBandwidthThreshold() > float64(1.1) {
			aggregationEpisode(bandwidth*2, time.Duration(50*time.Millisecond), 1000, true)
			Expect(tracker.numAckAggregationEpochs).To(Equal(uint64(3)))
		} else {
			aggregationEpisode(bandwidth*2, time.Duration(50*time.Millisecond), 1000, false)
			Expect(tracker.numAckAggregationEpochs).To(Equal(uint64(2)))
		}
	})

	It("SomewhatAggregatedSmallAcks", func() {
		aggregationEpisode(bandwidth*2, time.Duration(50*time.Millisecond), 100, true)
		aggregationEpisode(bandwidth*2, time.Duration(50*time.Millisecond), 100, true)
		now.Add(-1 * time.Millisecond)

		if tracker.AckAggregationBandwidthThreshold() > float64(1.1) {
			aggregationEpisode(bandwidth*2, time.Duration(50*time.Millisecond), 100, true)
			Expect(tracker.numAckAggregationEpochs).To(Equal(uint64(3)))
		} else {
			aggregationEpisode(bandwidth*2, time.Duration(50*time.Millisecond), 100, false)
			Expect(tracker.numAckAggregationEpochs).To(Equal(uint64(2)))
		}
	})

	It("NotAggregated", func() {
		aggregationEpisode(bandwidth, time.Duration(100*time.Millisecond), 100, true)
		Expect(tracker.numAckAggregationEpochs).To(Equal(uint64(2)))
	})

	It("StartNewEpochAfterAFullRound", func() {
		lastSentPacketNumber = protocol.PacketNumber(10)
		aggregationEpisode(bandwidth*2, time.Duration(50*time.Millisecond), 100, true)

		lastAckedPacketNumber = protocol.PacketNumber(11)
		// Update with a tiny bandwidth causes a very low expected bytes acked, which
		// in turn causes the current epoch to continue if the |tracker_| doesn't
		// check the packet numbers.
		tracker.Update(bandwidth/10, true, getRoundTripCount(), lastSentPacketNumber, lastAckedPacketNumber, now, 100)
		Expect(tracker.numAckAggregationEpochs).To(Equal(uint64(2)))
	})
})

var _ = Describe("BandwidthSampler", func() {
	var ()

	BeforeEach(func() {

	})

	It("test", func() {
		Expect(0).To(Equal(0))
		Expect(true).To(BeTrue())
		Expect(func() { panic("") }).To(Panic())

	})

	It("SendAndWait", func() {

	})

	It("SendTimeState", func() {

	})

	It("SendPaced", func() {

	})

	It("SendWithLosses", func() {

	})

	It("NotCongestionControlled", func() {

	})

	It("CompressedAck", func() {

	})

	It("ReorderedAck", func() {

	})

	It("AppLimited", func() {

	})

	It("FirstRoundTrip", func() {

	})

	It("RemoveObsoletePackets", func() {

	})

	It("NeuterPacket", func() {

	})

	It("CongestionEventSampleDefaultValues", func() {

	})

	It("TwoAckedPacketsPerEvent", func() {

	})

	It("LoseEveryOtherPacket", func() {

	})

	It("AckHeightRespectBandwidthEstimateUpperBound", func() {

	})
})
