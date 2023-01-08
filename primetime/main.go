package main

import (
	"bufio"
	"encoding/json"
	"errors"
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
			sc := bufio.NewScanner(c)
			for sc.Scan() {
				cdata := sc.Bytes()
				req, err := buildRequest(cdata)
				if err != nil {
					c.Write([]byte(fmt.Sprintf("invalid request data: %s\n", err)))
					c.Close()
					return
				}

				respData := buildResponse(req.Number)
				if err != nil {
					log.Println(err)
					c.Close()
					return
				}
				c.Write(respData)
			}
		}(conn)
	}
}

func isPrime(n int) bool {
	return big.NewInt(int64(n)).ProbablyPrime(0)
}

func buildRequest(data []byte) (request, error) {
	var req request
	if err := json.Unmarshal(data, &req); err != nil {
		return req, err
	}

	if req.Method != "isPrime" {
		return req, errors.New("invalid method")
	}

	return req, nil
}

func buildResponse(n int) []byte {
	resp := fmt.Sprintf(`{"method":"isPrime","prime":%t}`+"\n", isPrime(n))
	return []byte(resp)
}
