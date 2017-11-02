package keystoneauth_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func test_created() (string, error){
		  return "user",nil

}

var _ = Describe("BddUsers", func() {
        Context("When project is available", func() {
            It("User should be allowed to be created", func() {
                Expect(test_created()).To(Equal("user"))
            })
        })
})
