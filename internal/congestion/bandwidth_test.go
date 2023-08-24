package congestion

import (
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/quic-go/quic-go/internal/protocol"
)

var _ = Describe("Bandwidth", func() {
	It("converts from time delta", func() {
		Expect(BandwidthFromDelta(1, time.Millisecond)).To(Equal(1000 * BytesPerSecond))
	})

	It("bytes from bandwidth and time delta", func() {
		Expect(
			BytesFromBandwidthAndTimeDelta(
				Bandwidth(80000),
				100*time.Millisecond,
			)).To(Equal(protocol.ByteCount(1000)))
	})
})
