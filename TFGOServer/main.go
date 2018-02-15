package main

// main.go: server setup code

import (
	"net"
	"fmt"
	"os"
	"flag"
)

const (
	HOST = "128.135.175.185"
	PORT = ":9265"
)

// the following server code is heavily inspired by the example at
// https://astaxie.gitbooks.io/build-web-application-with-golang/en/08.1.html
func main() {
	flag.BoolVar(&verbose, "v", false, "print sent and received JSON")
	flag.Parse()

	tcpAddr, err := net.ResolveTCPAddr("tcp", HOST+PORT)
	checkError(err)
	listener, err := net.ListenTCP("tcp", tcpAddr)
	checkError(err)
	defer listener.Close()

	fmt.Println("Welcome to TFGO! Server is now awaiting connections.")

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
