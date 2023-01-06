package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"math/big"
	"net"
)

type request struct {
	Method string `json:"method"`
	Number int    `json:"number"`
}

type response struct {
	Method string `json:"method"`
	Prime  bool   `json:"prime"`
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
			defer c.Close()

			reader := bufio.NewReader(c)
			for {
				cdata, err := reader.ReadBytes('\n')
				if err != nil {
					log.Println(err)
					break
				}

				var req request
				if err := json.Unmarshal(cdata, &req); err != nil {
					log.Println(err)
					break
				}

				sdata, err := json.Marshal(response{Method: "isPrime", Prime: isPrime(req.Number)})
				if err != nil {
					log.Println(err)
					break
				}
				c.Write(sdata)
			}
		}(conn)
	}
}

func isPrime(n int) bool {
	return big.NewInt(int64(n)).ProbablyPrime(0)
}
