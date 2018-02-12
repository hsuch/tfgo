package main

import (
	"net"
	"fmt"
	"os"
)

const (
	HOST = "127.0.0.1"
	PORT = ":8080"
)

var games = make(map[string]*Game)

// the following server code is heavily inspired by the example at
// https://astaxie.gitbooks.io/build-web-application-with-golang/en/08.1.html
func main() {
	tcpAddr, err := net.ResolveTCPAddr("tcp", HOST+PORT)
	checkError(err)
	listener, err := net.ListenTCP("tcp", tcpAddr)
	checkError(err)
	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}
		go serveClient(conn)
	}
}

func checkError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}
}
