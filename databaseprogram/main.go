package main

import (
	"flag"
	"fmt"
	"log"
	"net"

	"github.com/mauricioabreu/protohackers/databaseprogram/protocol"
	"github.com/mauricioabreu/protohackers/databaseprogram/store"
)

func main() {
	port := flag.Int("port", 30001, "Specify the port run the server")
	flag.Parse()

	addr, err := net.ResolveUDPAddr("udp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatal(err)
	}

	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		log.Fatal(err)
	}

	defer conn.Close()

	log.Printf("database server running %s\n", conn.LocalAddr().String())

	kv := store.NewKeyValue("v1")

	for {
		buf := make([]byte, 1000)
		n, raddr, _ := conn.ReadFromUDP(buf)
		payload := protocol.ParseRequest(string(buf[0:n]))
		switch payload.Command {
		case protocol.Insert:
			kv.Insert(payload)
		case protocol.Retrieve:
			response := kv.Retrieve(payload.Key)
			conn.WriteToUDP([]byte(response), raddr)
		}
	}
}
