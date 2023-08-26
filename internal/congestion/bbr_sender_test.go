package congestion

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

//////////////////////////////////////////////////////////////////////
//                         quiche
//
//////////////////////////////////////////////////////////////////////

var _ = Describe("", func() {
	var ()

	BeforeEach(func() {
		Expect(true).To(BeTrue())
	})

	It("SetInitialCongestionWindow", func() {

	})

	It("SimpleTransfer", func() {

	})

	It("SimpleTransferBBRB", func() {

	})

	It("SimpleTransferSmallBuffer", func() {

	})

	It("RemoveBytesLostInRecovery", func() {

	})
	It("SimpleTransfer2RTTAggregationBytes", func() {

	})

	It("SimpleTransferAckDecimation", func() {

	})

	It("SimpleTransfer2RTTAggregationBytes20RTTWindow", func() {

	})

	It("SimpleTransfer2RTTAggregationBytes40RTTWindow", func() {

	})

	It("PacketLossOnSmallBufferStartup", func() {

	})

	It("PacketLossOnSmallBufferStartupDerivedCWNDGain", func() {

	})

	It("ApplicationLimitedBursts", func() {

	})

	It("ApplicationLimitedBurstsWithoutPrior", func() {

	})

	It("Drain", func() {

	})

	It("DISABLED_ShallowDrain", func() {

	})

	It("ProbeRtt", func() {

	})

	It("InFlightAwareGainCycling", func() {

	})

	It("SimpleTransfer1RTTStartup", func() {

	})

	It("SimpleTransfer2RTTStartup", func() {

	})

	It("SimpleTransferExitStartupOnLoss", func() {

	})

	It("SimpleTransferExitStartupOnLossSmallBuffer", func() {

	})

	It("DerivedPacingGainStartup", func() {

	})

	It("DerivedCWNDGainStartup", func() {

	})

	It("AckAggregationInStartup", func() {

	})

	It("ResumeConnectionState", func() {

	})

	It("ProbeRTTMinCWND1", func() {

	})

	It("StartupStats", func() {

	})

	It("RecalculatePacingRateOnCwndChange1RTT", func() {

	})

	It("RecalculatePacingRateOnCwndChange0RTT", func() {

	})

	It("MitigateCwndBootstrappingOvershoot", func() {

	})

	It("200InitialCongestionWindowWithNetworkParameterAdjusted", func() {

	})

	It("100InitialCongestionWindowWithNetworkParameterAdjusted", func() {

	})

	It("EnableOvershootingDetection", func() {

	})
})

//////////////////////////////////////////////////////////////////////
//                      quic-go cubic
//
//////////////////////////////////////////////////////////////////////

var _ = Describe("", func() {
	var ()

	BeforeEach(func() {

	})

	It("has the right values at startup", func() {

	})

	It("paces", func() {

	})

	It("application limited slow start", func() {

	})

	It("exponential slow start", func() {

	})

	It("slow start packet loss", func() {

	})

	It("slow start packet loss PRR", func() {

	})

	It("slow start burst packet loss PRR", func() {

	})

	It("RTO congestion window", func() {

	})

	It("RTO congestion window no retransmission", func() {

	})

	It("tcp cubic reset epoch on quiescence", func() {

	})

	It("multiple losses in one window", func() {

	})

	It("1 connection congestion avoidance at end of recovery", func() {

	})

	It("no PRR", func() {

	})

	It("reset after connection migration", func() {

	})

	It("slow starts up to the maximum congestion window", func() {

	})

	It("doesn't allow reductions of the maximum packet size", func() {

	})

	It("limit cwnd increase in congestion avoidance", func() {

	})
})

//////////////////////////////////////////////////////////////////////
//                      mvfst
//
//////////////////////////////////////////////////////////////////////

var _ = Describe("", func() {
	var ()

	BeforeEach(func() {

	})

	It("InitStates", func() {

	})

	It("InitWithCwndAndRtt", func() {

	})

	It("Recovery", func() {

	})

	It("StartupCwnd", func() {

	})

	It("StartupCwndImplicit", func() {

	})

	It("LeaveStartup", func() {

	})

	It("RemoveInflightBytes", func() {

	})

	It("ProbeRtt", func() {

	})

	It("NoLargestAckedPacketInitialNoCrash", func() {

	})

	It("NoLargestAckedPacketHandshakeNoCrash", func() {

	})

	It("NoLargestAckedPacketAppDataNoCrash", func() {

	})

	It("NoLargestAckedPacketInitialNoCrashPn1", func() {

	})

	It("NoLargestAckedPacketHandshakeNoCrashPn1", func() {

	})

	It("NoLargestAckedPacketAppDataNoCrashPn1", func() {

	})

	It("AckAggregation", func() {

	})

	It("AppLimited", func() {

	})

	It("AppLimitedIgnored", func() {

	})

	It("ExtendMinRttExpiration", func() {

	})

	It("BytesCounting", func() {

	})

	It("AppIdle", func() {

	})

	It("PacketLossInvokesPacer", func() {

	})

	It("ProbeRttSetsAppLimited", func() {

	})

	It("BackgroundMode", func() {

	})

	It("GetBandwidthSample", func() {

	})
})
