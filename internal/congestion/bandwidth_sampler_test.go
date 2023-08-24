package congestion

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

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

var _ = Describe("Max ack height tracker", func() {
	var ()

	BeforeEach(func() {

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
