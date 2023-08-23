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

		initializeMinFilter = func() {
			var nowTime int64 = 0
			var rttSample time.Duration = 10
			for i := 0; i < 5; i++ {
				windowedMinRtt.Update(rttSample, nowTime)
				nowTime += 25
				rttSample += 10
			}
		}

		initializeMaxFilter = func() {
			var nowTime int64 = 0
			var bwSample uint64 = 10
			for i := 0; i < 5; i++ {
				windowedMaxBw.Update(bwSample, nowTime)
				nowTime += 25
				bwSample -= 100
			}
		}
	)

	BeforeEach(func() {
		// now = time.Now().UnixNano() / 1e6
		windowedMinRtt = NewWindowedFilter[time.Duration, int64](10, MinFilter[time.Duration])
		windowedMaxBw = NewWindowedFilter[uint64, int64](1000, MaxFilter[uint64])
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
		initializeMinFilter()
		initializeMaxFilter()
	})
})
