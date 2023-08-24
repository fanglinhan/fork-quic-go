package congestion

import (
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/quic-go/quic-go/internal/protocol"
)

var _ = Describe("Bandwidth sampler", func() {
	var (
		//now       time.Time
		bandwidth Bandwidth = Bandwidth(10 * 1000)

		aggregationEpisode = func(
			aggregationBandwidth Bandwidth,
			aggregationDuration time.Duration,
			bytesPerAck protocol.ByteCount,
			expectNewAggregationEpoch bool,
		) {
			Expect(aggregationBandwidth >= bandwidth).To(BeTrue())
			//startTime := now
			// aggregationBytes := aggregationBandwidth * aggregationDuration
		}
	)

	BeforeEach(func() {

	})

	It("test", func() {
		Expect(0).To(Equal(0))
		Expect(true).To(BeTrue())
		Expect(func() { panic("") }).To(Panic())

	})

	It("SendAndWait", func() {
		aggregationEpisode(bandwidth, time.Duration(6*time.Millisecond), 1200, true)
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

var _ = Describe("Max ack height tracker", func() {
	var (
		tracker *maxAckHeightTracker
	)

	BeforeEach(func() {
		tracker = newMaxAckHeightTracker(10)
		tracker.SetAckAggregationBandwidthThreshold(float64(1.8))
		tracker.SetStartNewAggregationEpochAfterFullRound(true)
	})

	It("VeryAggregatedLargeAck", func() {

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
