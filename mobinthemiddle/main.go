package main

import (
	"flag"
	"fmt"
	"log"
	"net"

	"github.com/mauricioabreu/protohackers/mobinthemiddle/proxy"
)

func main() {
	port := flag.Int("port", 30001, "Specify the port run the server")
	chatAddress := flag.String("chat_address", "chat.protohackers.com:16963", "Address of upstream chat server")
	flag.Parse()

	l, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatal(err)
	}
	defer l.Close()

	for {
		rc, err := net.Dial("tcp", *chatAddress)
		if err != nil {
			log.Printf("Error connecting to %s: %s", *chatAddress, err.Error())
			return
		}

		lc, err := l.Accept()
		if err != nil {
			log.Printf("Error accepting connection: %s", err.Error())
			return
		}

		p := proxy.New(lc, rc)
		go p.Handle()
	}
}
