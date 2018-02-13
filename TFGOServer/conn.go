package main

import (
	"net"
	"encoding/json"
)

type ClientMessage struct {
	Action string
	Data interface{}
}

func serveClient(conn net.Conn) {
	defer conn.Close()

	d := json.NewDecoder(conn)
	for {
		var msg ClientMessage
		if err := d.Decode(&msg); err != nil {
			break
		}


	}
}
