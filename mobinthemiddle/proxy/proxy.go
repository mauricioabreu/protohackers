package proxy

import (
	"bufio"
	"log"
	"net"
	"regexp"
	"strings"
)

var (
	RgxBogusCoin          = regexp.MustCompile(`^7[A-Za-z0-9]{25,34}$`)
	MaliciousBogusAddress = "7YWHMfk9JZe0LM0g1ZauHuiSxhI"
)

type Proxy struct {
	LocalConn  net.Conn
	RemoteConn net.Conn
}

func New(lc, rc net.Conn) *Proxy {
	return &Proxy{
		LocalConn:  lc,
		RemoteConn: rc,
	}
}

func (p *Proxy) RewriteBogusCoin(message, replace string) string {
	rewritten := false
	words := strings.Split(strings.Trim(message, "\n"), " ")
	for k, word := range words {
		if RgxBogusCoin.MatchString(word) {
			words[k] = replace
			rewritten = true
		}
	}
	if rewritten {
		return strings.Join(words, " ") + "\n"
	}
	return message
}

func (p *Proxy) Handle() {
	go p.HandleMessage(p.LocalConn, p.RemoteConn)
	p.HandleMessage(p.RemoteConn, p.LocalConn)
}

func (p *Proxy) HandleMessage(from, to net.Conn) {
	reader := bufio.NewReader(from)

	defer to.Close()

	for {
		in, err := reader.ReadString('\n')
		if err != nil {
			log.Println(err)
			return
		}

		msg := p.RewriteBogusCoin(in, MaliciousBogusAddress)
		if _, err := to.Write([]byte(msg)); err != nil {
			log.Printf("Failed to write message: %s\n", err.Error())
		}
	}
}
