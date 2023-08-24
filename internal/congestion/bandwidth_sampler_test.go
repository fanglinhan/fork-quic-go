package congestion

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("unit", func() {
	var ()

	BeforeEach(func() {

	})

	It("test", func() {
		Expect(0).To(Equal(0))
		Expect(true).To(BeTrue())
		Expect(func() { panic("") }).To(Panic())
	})
})
