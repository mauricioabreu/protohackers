package proxy_test

import (
	"net"
	"testing"

	"github.com/mauricioabreu/protohackers/mobinthemiddle/proxy"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestProxy(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Proxy Suite")
}

func newConn() net.Conn {
	_, s := net.Pipe()
	return s
}

var _ = Describe("Proxy", func() {
	Describe("Rewrite BogusCoin message", func() {
		When("message does not have a BogusCoin address", func() {
			It("keeps the message untouched", func() {
				p := proxy.New(newConn(), newConn())
				message := "Hi alice, please send payment to 1234foobarZzZz"
				Expect(p.RewriteBogusCoin(message, "7YWHMfk9JZe0LM0g1ZauHuiSxhI")).To(Equal(message))
			})
		})
		When("message has a BogusCoin address", func() {
			It("rewrites the address", func() {
				p := proxy.New(newConn(), newConn())
				originalMessage := "Hi alice, please send payment to 7iKDZEwPZSqIvDnHvVN2r0hUWXD5rHX"
				maliciousMessage := "Hi alice, please send payment to 7YWHMfk9JZe0LM0g1ZauHuiSxhI\n"
				Expect(p.RewriteBogusCoin(originalMessage, "7YWHMfk9JZe0LM0g1ZauHuiSxhI")).To(Equal(maliciousMessage))
			})
		})
		When("message has multiple BogusCoin address", func() {
			It("rewrites all the address", func() {
				p := proxy.New(newConn(), newConn())
				originalMessage := "Please pay the ticket price of 15 Boguscoins to one of these addresses: 7YWHMfk9JZe0LM0g1ZauHuiSxhI 74wRkQ2E0lg1v2j779Cf8cHmK6"
				maliciousMessage := "Please pay the ticket price of 15 Boguscoins to one of these addresses: 7YWHMfk9JZe0LM0g1ZauHuiSxhI 7YWHMfk9JZe0LM0g1ZauHuiSxhI\n"
				Expect(p.RewriteBogusCoin(originalMessage, "7YWHMfk9JZe0LM0g1ZauHuiSxhI")).To(Equal(maliciousMessage))
			})
		})
	})
})
