package store

import (
	"fmt"

	"github.com/mauricioabreu/protohackers/databaseprogram/protocol"
)

type Store interface {
	Insert()
	Retrieve() string
	Version() string
}

type KeyValue struct {
	data    map[string]string
	version string
}

func NewKeyValue(version string) *KeyValue {
	return &KeyValue{
		data:    make(map[string]string),
		version: version,
	}
}

func (kv *KeyValue) Insert(p protocol.Payload) {
	kv.data[p.Key] = p.Value
}

func (kv *KeyValue) Retrieve(key string) string {
	if key == "version" {
		return kv.Version()
	}
	value := kv.data[key]
	return fmt.Sprintf("%s=%s", key, value)
}

func (kv *KeyValue) Version() string {
	return fmt.Sprintf("version=%s", kv.version)
}
