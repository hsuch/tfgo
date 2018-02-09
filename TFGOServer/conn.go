package TFGOServer

import (
	"net"
)

func serveClient(conn net.Conn) {
	defer conn.Close()
}
