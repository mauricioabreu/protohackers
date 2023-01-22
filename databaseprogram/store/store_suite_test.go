package store_test

import (
	"testing"

	"github.com/mauricioabreu/protohackers/databaseprogram/protocol"
	"github.com/mauricioabreu/protohackers/databaseprogram/store"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestStore(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Store Suite")
}

var _ = Describe("Store", func() {
	Describe("Key value database", func() {
		When("insert is called", func() {
			It("inserts the value in the given key", func() {
				kvd := store.NewKeyValue("v1")
				kvd.Insert(protocol.Payload{Key: "foo", Value: "bar"})
			})
		})
		When("retrieve is called with a regular key", func() {
			It("retrieves the value for the given key", func() {
				kvd := store.NewKeyValue("v1")
				kvd.Insert(protocol.Payload{Key: "foo", Value: "bar"})
				Expect(kvd.Retrieve("foo")).To(Equal("foo=bar"))
			})
		})
		When("retrieve is called with a version argument", func() {
			It("retrieves the database version", func() {
				kvd := store.NewKeyValue("v1")
				Expect(kvd.Retrieve("version")).To(Equal("version=v1"))
			})
		})
		When("client wants to know the version", func() {
			It("retrieves the database version", func() {
				kvd := store.NewKeyValue("v1")
				Expect(kvd.Version()).To(Equal("version=v1"))
			})
		})
	})
})
