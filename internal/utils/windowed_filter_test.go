package utils

import (
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("unit", func() {
	var (
		// now            int64
		windowedMinRtt *WindowedFilter[time.Duration, int64]
		windowedMaxBw  *WindowedFilter[uint64, int64]

		// initializeMinFilter = func() {
		// 	var nowTime int64 = 0
		// 	var rttSample time.Duration = 10
		// 	for i := 0; i < 5; i++ {
		// 		windowedMinRtt.Update(rttSample, nowTime)
		// 		nowTime += 25
		// 		rttSample += 10
		// 	}
		// }

		// initializeMaxFilter = func() {
		// 	var nowTime int64 = 0
		// 	var bwSample uint64 = 10
		// 	for i := 0; i < 5; i++ {
		// 		windowedMaxBw.Update(bwSample, nowTime)
		// 		nowTime += 25
		// 		bwSample -= 100
		// 	}
		// }

		// updateWithIrrelevantSamples = func() {

		// }

	)

	BeforeEach(func() {
		// now = time.Now().UnixNano() / 1e6
		windowedMinRtt = NewWindowedFilter[time.Duration, int64](99, MinFilter[time.Duration])
		windowedMaxBw = NewWindowedFilter[uint64, int64](99, MaxFilter[uint64])
	})

	It("UninitializedEstimates", func() {
		Expect(windowedMinRtt.GetBest()).To(Equal(time.Duration(0)))
		Expect(windowedMinRtt.GetSecondBest()).To(Equal(time.Duration(0)))
		Expect(windowedMinRtt.GetThirdBest()).To(Equal(time.Duration(0)))
		Expect(windowedMaxBw.GetBest()).To(Equal(uint64(0)))
		Expect(windowedMaxBw.GetSecondBest()).To(Equal(uint64(0)))
		Expect(windowedMaxBw.GetThirdBest()).To(Equal(uint64(0)))
	})

	It("MonotonicallyIncreasingMin", func() {
		var nowTime int64 = 0
		var rttSample time.Duration = 10
		windowedMinRtt.Update(rttSample, nowTime)
		Expect(windowedMinRtt.GetBest()).To(Equal(time.Duration(10)))

		// Gradually increase the rtt samples and ensure the windowed min rtt starts
		// rising.
		for i := 0; i < 6; i++ {
			nowTime += 25
			rttSample += 10
			windowedMinRtt.Update(rttSample, nowTime)
			if i < 3 {
				Expect(windowedMinRtt.GetBest()).To(Equal(time.Duration(10)))
			} else if i == 3 {
				Expect(windowedMinRtt.GetBest()).To(Equal(time.Duration(20)))
			} else if i < 6 {
				Expect(windowedMinRtt.GetBest()).To(Equal(time.Duration(40)))
			}
		}
	})

})
