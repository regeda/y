package y

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Collection", func() {
	type something struct {
		ID int64 `y:"id,pk"`
	}

	Context("when one item added", func() {
		var (
			ptr *something
			c   *Collection
		)

		BeforeEach(func() {
			ptr = &something{ID: 1}
			c = New(ptr).Collection()
		})

		It("should be non-empty", func() {
			Expect(c.Empty()).To(BeFalse())
		})

		It("should be contain the first item", func() {
			Expect(c.First()).To(Equal(ptr))
		})
	})

	Context("when two items added", func() {
		var (
			ptrs []*something
			c    *Collection
		)

		BeforeEach(func() {
			ptrs = make([]*something, 2)
			for i, id := range []int64{1, 2} {
				ptr := &something{ID: id}
				ptrs[i] = ptr
			}
			c = New(ptrs).Collection()
		})

		It("should contain two items", func() {
			Expect(c.Size()).To(Equal(2))
		})

		It("should list correct sequence of items", func() {
			Expect(c.List()).To(Equal(ptrs))
		})
	})
})
