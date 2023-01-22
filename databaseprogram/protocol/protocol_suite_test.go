package protocol_test

import (
	"testing"

	"github.com/mauricioabreu/protohackers/databaseprogram/protocol"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestProtocol(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Protocol Suite")
}

var _ = Describe("Protocol", func() {
	Describe("Parsing insert request", func() {
		When("payload has only one equals sign", func() {
			It("builds a correct key value", func() {
				p := protocol.ParseRequest("foo=bar")
				Expect(p.Key).To(Equal("foo"))
				Expect(p.Value).To(Equal("bar"))
			})
		})
		When("payload has two distant equal signs", func() {
			It("builds a correct key value", func() {
				p := protocol.ParseRequest("foo=bar=baz")
				Expect(p.Key).To(Equal("foo"))
				Expect(p.Value).To(Equal("bar=baz"))
			})
		})
		When("payload has nothing after equals sign", func() {
			It("builds a correct key value", func() {
				p := protocol.ParseRequest("foo=")
				Expect(p.Key).To(Equal("foo"))
				Expect(p.Value).To(Equal(""))
			})
		})
		When("payload has multiple consecutive equal signs", func() {
			It("builds a correct key value", func() {
				p := protocol.ParseRequest("foo===")
				Expect(p.Key).To(Equal("foo"))
				Expect(p.Value).To(Equal("=="))
			})
		})
		When("payload has nothing before equal signs", func() {
			It("builds a correct key value", func() {
				p := protocol.ParseRequest("=foo")
				Expect(p.Key).To(Equal(""))
				Expect(p.Value).To(Equal("foo"))
			})
		})
	})
})
