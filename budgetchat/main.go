package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"net"
	"strings"

	"github.com/mauricioabreu/protohackers/budgetchat/broker"
)

func main() {
	port := flag.Int("port", 30001, "Specify the port run the server")
	flag.Parse()

	l, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatal(err)
	}
	defer l.Close()

	b := broker.New()

	for {
		conn, err := l.Accept()
		if err != nil {
			log.Fatal(err)
		}

		go handle(b, conn)
	}
}

func handle(b *broker.Broker, conn net.Conn) {
	fmt.Fprintln(conn, "Welcome to budgetchat! What shall I call you?")

	reader := bufio.NewReader(conn)
	username, err := reader.ReadString('\n')
	if err != nil {
		fmt.Fprintln(conn, err)
		conn.Close()
		return
	}
	username = strings.TrimSpace(username)

	if err := b.AddUser(username, conn); err != nil {
		log.Println(err)
		conn.Close()
		return
	}

	defer b.RemoveUser(username)

	for {
		in, err := reader.ReadString('\n')
		if err != nil {
			log.Println(err)
			conn.Close()
			return
		}
		b.SendMessage(username, strings.TrimSpace(in))
	}
}
