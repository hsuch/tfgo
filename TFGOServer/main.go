package main

// main.go: server setup code

import (
	"net"
	"fmt"
	"os"
	"flag"
)

// used to suppress sends to clients while testing
var isTesting = false

// the following server code is heavily inspired by the example at
// https://astaxie.gitbooks.io/build-web-application-with-golang/en/08.1.html
func main() {
	flag.IntVar(&verbosity, "v", 0, "verbosity: 0 (high-level info only), 1 (non-periodic JSON), or 2 (all JSON)")
	var HOST = flag.String("host", "127.0.0.1", "server ip")
	var PORT = flag.String("port", "9265", "server port")
	flag.Parse()

	tcpAddr, err := net.ResolveTCPAddr("tcp", *HOST+":"+*PORT)
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
