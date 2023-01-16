package broker

import (
	"errors"
	"fmt"
	"io"
	"log"
	"regexp"
	"sort"
	"sync"

	"golang.org/x/exp/maps"
)

var (
	RgxUsername           = regexp.MustCompile(`^[a-zA-Z0-9]{1,}$`)
	ErrDuplicatedUsername = errors.New("duplicated user name")
	ErrInvalidUsername    = errors.New("invalid username")
)

type Broker struct {
	sync.Mutex
	users map[string]io.Writer
}

func New() *Broker {
	return &Broker{users: make(map[string]io.Writer)}
}

func (b *Broker) AddUser(username string, conn io.Writer) error {
	if _, ok := b.users[username]; ok {
		return ErrDuplicatedUsername
	}

	if err := validateUserName(username); err != nil {
		return err
	}

	b.Lock()
	b.users[username] = conn
	b.Unlock()

	b.Broadcast(username, fmt.Sprintf("* %s has entered the room\n", username))
	fmt.Fprintf(conn, "* the room contains: %s\n", b.ListUsers(username))

	return nil
}

func (b *Broker) RemoveUser(username string) {
	b.Lock()
	delete(b.users, username)
	b.Unlock()
	b.Broadcast(username, fmt.Sprintf("* %s has left the room\n", username))
}

func (b *Broker) Broadcast(sender, message string) {
	for user, conn := range b.users {
		if sender != user {
			if _, err := fmt.Fprint(conn, message); err != nil {
				log.Printf("Failed to broadcast message to other users: %s", err)
			}
		}
	}
}

func (b *Broker) ListUsers(except string) []string {
	users := maps.Clone(b.users)
	delete(users, except)

	names := maps.Keys(users)

	sort.Strings(names)
	return names
}

func (b *Broker) SendMessage(username, message string) {
	b.Broadcast(username, fmt.Sprintf("[%s] %s\n", username, message))
}

func validateUserName(username string) error {
	if !RgxUsername.MatchString(username) {
		return ErrInvalidUsername
	}
	return nil
}
