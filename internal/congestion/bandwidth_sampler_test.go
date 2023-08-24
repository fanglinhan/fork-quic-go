package congestion

import (
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/quic-go/quic-go/internal/protocol"
)

var _ = Describe("function", func() {
	It("bytes from bandwidth and time delta", func() {
		Expect(
			BytesFromBandwidthAndTimeDelta(
				Bandwidth(80000),
				100*time.Millisecond,
			)).To(Equal(protocol.ByteCount(1000)))
	})

	It("time delta from bytes and bandwidth", func() {
		Expect(TimeDeltaFromBytesAndBandwidth(
			protocol.ByteCount(50000),
			Bandwidth(400),
		)).To(Equal(1000 * time.Second))
	})
})

var _ = Describe("Max ack height tracker", func() {
	var (
		tracker *maxAckHeightTracker
		//now       time.Time
		bandwidth Bandwidth = Bandwidth(10 * 1000)

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
			// Expect(aggregationBandwidth >= bandwidth).To(BeTrue())
			// startTime := now
			// aggregationBytes := BytesFromBandwidthAndTimeDelta(aggregationBandwidth, aggregationDuration)
			// numAcks := aggregationBytes / bytesPerAck
			// Expect(aggregationBytes).To(Equal(numAcks * bytesPerAck))
			// timeBetweenAcks := aggregationDuration / time.Duration(numAcks)
			// Expect(aggregationDuration).To(Equal(time.Duration(numAcks) * timeBetweenAcks))

			// The total duration of aggregation time and quiet period.
			//totalDuration =

		}
	)

	BeforeEach(func() {
		tracker = newMaxAckHeightTracker(10)
		tracker.SetAckAggregationBandwidthThreshold(float64(1.8))
		tracker.SetStartNewAggregationEpochAfterFullRound(true)
	})

	It("VeryAggregatedLargeAck", func() {
		aggregationEpisode(bandwidth, time.Duration(6*time.Millisecond), 1200, true)
	})

	It("VeryAggregatedSmallAcks", func() {

	})

	It("SomewhatAggregatedLargeAck", func() {

	})

	It("SomewhatAggregatedSmallAcks", func() {

	})

	It("NotAggregated", func() {

	})

	It("StartNewEpochAfterAFullRound", func() {

	})
})

var _ = Describe("Bandwidth sampler", func() {
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
