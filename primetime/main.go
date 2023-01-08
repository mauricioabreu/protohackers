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
	Method string   `json:"method"`
	Number *float64 `json:"number"`
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

				respData := buildResponse(*req.Number)
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

func buildRequest(data []byte) (request, error) {
	var req request
	if err := json.Unmarshal(data, &req); err != nil {
		return req, err
	}

	if req.Method != "isPrime" {
		return req, errors.New("invalid method")
	}

	if req.Number == nil {
		return req, errors.New("missing number")
	}

	return req, nil
}

func buildResponse(n float64) []byte {
	resp := fmt.Sprintf(`{"method":"isPrime","prime":%t}`+"\n", isPrime(n))
	return []byte(resp)
}

func isPrime(n float64) bool {
	return big.NewInt(int64(n)).ProbablyPrime(0)
}
