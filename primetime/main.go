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
			reader := bufio.NewReader(c)
			cdata, err := reader.ReadBytes('\n')
			if err != nil {
				log.Println(err)
				return
			}

			req, err := buildRequest(cdata)
			if err != nil {
				c.Write([]byte(fmt.Sprintf("invalid request data: %s", err)))
				c.Close()
				return
			}

			respData, err := buildResponse(req.Number)
			if err != nil {
				log.Println(err)
				c.Close()
				return
			}
			c.Write(respData)
		}(conn)
	}
}

func isPrime(n int) bool {
	return big.NewInt(int64(n)).ProbablyPrime(0)
}

func buildResponse(n int) ([]byte, error) {
	data, err := json.Marshal(response{Method: "isPrime", Prime: isPrime(n)})
	if err != nil {
		return data, err
	}
	return data, nil
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
