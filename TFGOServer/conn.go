package main

import (
	"net"
)

func serveClient(conn net.Conn) {
	defer conn.Close()
}
