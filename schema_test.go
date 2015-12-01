package y

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Schema", func() {
	var p *Proxy

	Context("when a struct has one field without tags", func() {
		type something struct {
			ID int64
		}

		BeforeEach(func() {
			p = New(something{})
		})

		It("should contain one parsed field", func() {
			Expect(p.schema.fields).To(HaveLen(1))
		})

		It("should contain not nil field schema", func() {
			Expect(p.schema.fields["id"]).NotTo(BeNil())
		})
	})

	Context("when a struct has disabled field", func() {
		type something struct {
			ID int64 `y:"-"`
		}

		BeforeEach(func() {
			p = New(something{})
		})

		It("should be empty", func() {
			Expect(p.schema.fields).To(BeEmpty())
		})
	})

	Context("when a struct has a primary key without name", func() {
		type something struct {
			ID int64 `y:",pk"`
		}

		BeforeEach(func() {
			p = New(something{})
		})

		It("should contain a primary key with name", func() {
			Expect(p.schema.xinfo.pk).To(Equal([]string{"id"}))
		})
	})

	Context("when a struct has a composite primary key", func() {
		type something struct {
			X int64 `y:"x,pk"`
			Y int64 `y:"y,pk"`
		}

		BeforeEach(func() {
			p = New(something{})
		})

		It("should contain correct PK sequence", func() {
			Expect(p.schema.xinfo.pk).To(Equal([]string{"x", "y"}))
		})

		It("should contain the first key in the index", func() {
			Expect(p.schema.xinfo.idx["x"]).To(Equal(1))
		})

		It("shouldn't contain the second key in the index", func() {
			Expect(p.schema.xinfo.idx["y"]).To(Equal(0))
		})
	})

	Context("when a struct has an auto-incremented key", func() {
		type something struct {
			ID int64 `y:",autoincr"`
		}

		BeforeEach(func() {
			p = New(something{})
		})

		It("should contains enabled autoincr in field opts", func() {
			Expect(p.schema.fields["id"].autoincr).To(BeTrue())
		})
	})

	Context("when a slice of struct parsed", func() {
		type something struct {
			ID int64
		}

		BeforeEach(func() {
			p = New([]something{})
		})

		It("should contain one parsed field", func() {
			Expect(p.schema.fields).To(HaveLen(1))
		})

		It("should contain not nil field schema", func() {
			Expect(p.schema.fields["id"]).NotTo(BeNil())
		})
	})

	Context("when a struct contain embedded struct", func() {
		type foo struct {
			x int64
			y int64
		}
		type bar struct {
			foo
			z int64
		}

		BeforeEach(func() {
			p = New(bar{})
		})

		It("should contain three parsed field", func() {
			Expect(p.schema.fields).To(HaveLen(3))
		})

		It("should contain parsed fields opts", func() {
			Expect(p.schema.fields["x"]).NotTo(BeNil())
			Expect(p.schema.fields["y"]).NotTo(BeNil())
			Expect(p.schema.fields["z"]).NotTo(BeNil())
		})
	})

	Context("when a struct has a foreign key", func() {
		type something struct {
			AnyID int64 `y:",fk"`
		}

		BeforeEach(func() {
			p = New(something{})
		})

		It("should contain index opts for 'any_id'", func() {
			Expect(p.schema.xinfo.idx["any_id"]).To(Equal(1))
		})

		It("should contain 'any_id' as own key", func() {
			Expect(p.schema.fks["any"].from).To(Equal("any_id"))
		})

		It("should contain 'id' as a target key", func() {
			Expect(p.schema.fks["any"].target).To(Equal("id"))
		})
	})

	Context("when a struct has overrided foreign key", func() {
		type something struct {
			AnyID int64 `y:"aid,fk:any.id"`
		}

		BeforeEach(func() {
			p = New(something{})
		})

		It("should contain index opts for 'aid'", func() {
			Expect(p.schema.xinfo.idx["aid"]).To(Equal(1))
		})

		It("should contain 'aid' as own key", func() {
			Expect(p.schema.fks["any"].from).To(Equal("aid"))
		})

		It("should contain 'id' as a target key", func() {
			Expect(p.schema.fks["any"].target).To(Equal("id"))
		})
	})
})
