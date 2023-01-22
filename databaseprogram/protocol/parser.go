package protocol

import "strings"

type Command int

const (
	Insert Command = iota + 1
	Retrieve
)

type Payload struct {
	Key     string
	Value   string
	Command Command
}

func ParseRequest(data string) Payload {
	key, value, found := strings.Cut(data, "=")
	payload := Payload{Key: key, Value: value}
	if found {
		payload.Command = Insert
	} else {
		payload.Command = Retrieve
	}
	return payload
}
