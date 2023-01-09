package main

import (
	"encoding/binary"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
)

type message struct {
	Type uint8
	Arg1 int32
	Arg2 int32
}

func main() {
	port := flag.Int("port", 30001, "Specify the port run the server")
	flag.Parse()

	l, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatal(err)
	}
	defer l.Close()

	for {
		conn, err := l.Accept()
		if err != nil {
			log.Fatal(err)
		}

		go func(c net.Conn) {
			prices := map[int32]int32{}
			buf := make([]byte, 9)

			for {
				if _, err := io.ReadAtLeast(c, buf, 9); err != nil {
					return
				}

				msg := decodePayload(buf)
				switch msg.Type {
				case 'I':
					prices[msg.Arg1] = msg.Arg2
				case 'Q':
					var sum, total int32

					for time, price := range prices {
						if time >= msg.Arg1 && time <= msg.Arg2 {
							total++
							sum += price
						}
					}

					if total == 0 {
						total = 1
					}

					mean := sum / total
					ans := make([]byte, 4)
					binary.BigEndian.PutUint32(ans, uint32(mean))
					c.Write(ans)
				}
			}
		}(conn)
	}
}

func decodePayload(data []byte) message {
	return message{
		Type: data[0],
		Arg1: int32(binary.BigEndian.Uint32(data[1:5])),
		Arg2: int32(binary.BigEndian.Uint32(data[5:])),
	}
}
