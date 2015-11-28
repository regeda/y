package y

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Collection", func() {
	type something struct {
		ID int64 `db:"id,pk"`
	}

	var (
		c *Collection
		p *Proxy
	)

	BeforeEach(func() {
		p = New(something{})
		c = p.blankCollection()
	})

	It("should be empty", func() {
		Expect(c.Empty()).To(BeTrue())
	})

	Context("when one item added", func() {
		var ptr interface{}

		BeforeEach(func() {
			v := p.schema.create()
			v.field("ID").SetInt(1)
			v.addTo(c)
			ptr = v.ptr().Interface()
		})

		It("should be non-empty", func() {
			Expect(c.Empty()).To(BeFalse())
		})

		It("should be contain the first item", func() {
			Expect(c.First()).To(Equal(ptr))
		})
	})

	Context("when two items added", func() {
		var ptrs []interface{}

		BeforeEach(func() {
			ptrs = []interface{}{}
			for _, id := range []int64{1, 2} {
				v := p.schema.create()
				v.field("ID").SetInt(id)
				v.addTo(c)
				ptrs = append(ptrs, v.ptr().Interface())
			}
		})

		It("should contain two items", func() {
			Expect(c.Size()).To(Equal(2))
		})

		It("should list correct sequence of items", func() {
			Expect(c.List()).To(Equal(ptrs))
		})
	})
})
