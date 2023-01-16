package broker_test

import (
	"bytes"
	"testing"

	"github.com/mauricioabreu/protohackers/budgetchat/broker"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestBroker(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Broker Suite")
}

var _ = Describe("Broker", func() {
	Describe("Adding a new user", func() {
		When("user already exists", func() {
			It("returns an error of duplicated user", func() {
				var buf bytes.Buffer
				b := broker.New()

				Expect(b.AddUser("john", &buf)).To(Not(HaveOccurred()))
				Expect(b.AddUser("john", &buf)).To(Equal(broker.ErrDuplicatedUsername))
			})
		})
		When("username is empty", func() {
			It("returns an error of invalid name", func() {
				var buf bytes.Buffer
				b := broker.New()

				Expect(b.AddUser("", &buf)).To(Equal(broker.ErrInvalidUsername))
			})
		})
		When("username contains invalid characters", func() {
			It("returns an error of invalid name", func() {
				var buf bytes.Buffer
				b := broker.New()

				Expect(b.AddUser("john?", &buf)).To(Equal(broker.ErrInvalidUsername))
			})
		})
		When("username contains multiple names separated by space", func() {
			It("returns an error of invalid name", func() {
				var buf bytes.Buffer
				b := broker.New()

				Expect(b.AddUser("john ana", &buf)).To(Equal(broker.ErrInvalidUsername))
			})
		})
		When("username is composed only by spaces", func() {
			It("returns an error of invalid name", func() {
				var buf bytes.Buffer
				b := broker.New()

				Expect(b.AddUser("    ", &buf)).To(Equal(broker.ErrInvalidUsername))
			})
		})
	})
	Describe("Listing users", func() {
		When("chat has no users", func() {
			It("returns empty list", func() {
				b := broker.New()

				Expect(len(b.ListUsers(""))).To(Equal(0))
			})
		})
		When("chat contains some users", func() {
			It("returns a list of names alphabetically sorted, except the sender", func() {
				var buf bytes.Buffer
				b := broker.New()
				b.AddUser("john", &buf)
				b.AddUser("kelly", &buf)
				b.AddUser("ana", &buf)

				Expect(b.ListUsers("ana")).To(Equal([]string{"john", "kelly"}))
			})
		})
	})
	Describe("Sending a message", func() {
		It("broadcast the message to all connected users", func() {
			var abuf, jbuf bytes.Buffer
			b := broker.New()
			b.AddUser("ana", &abuf)
			b.AddUser("john", &jbuf)

			b.SendMessage("ana", "hello guys, how you doing?")

			Expect(abuf.String()).To(Not(ContainSubstring("[ana] hello guys, how you doing?")))
			Expect(jbuf.String()).To(ContainSubstring("[ana] hello guys, how you doing?"))
		})
	})
	Describe("Leaving the chat", func() {
		It("broadcasts to others the user left", func() {
			var abuf, jbuf bytes.Buffer
			b := broker.New()
			b.AddUser("ana", &abuf)
			b.AddUser("john", &jbuf)

			b.RemoveUser("ana")

			Expect(abuf.String()).To(Not(ContainSubstring("* ana has left the room")))
			Expect(jbuf.String()).To(ContainSubstring("* ana has left the room"))
		})
	})
})
