package y

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Utils", func() {
	It("should be snake", func() {
		Expect(underscore("ILoveGoAndJSONSoMuch")).To(Equal("i_love_go_and_json_so_much"))
		Expect(underscore("CamelCase")).To(Equal("camel_case"))
		Expect(underscore("Camel")).To(Equal("camel"))
		Expect(underscore("CAMEL")).To(Equal("camel"))
		Expect(underscore("camel")).To(Equal("camel"))
		Expect(underscore("BIGCase")).To(Equal("big_case"))
		Expect(underscore("privateCase")).To(Equal("private_case"))
		Expect(underscore("PublicCASE")).To(Equal("public_case"))
	})
})
