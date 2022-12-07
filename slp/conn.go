package slp

import (
	"fmt"
	"net"
)

// handle a connection
func HandleConnection(conn net.Conn, getStatus func(net.Conn) StatusResponse, getMsg func(net.Conn, LoginStart) Disconnect) {
	defer conn.Close()
	// fmt.Println("--- New Connection ---")
	// defer fmt.Print("--- Connection Closed ---\n\n")

	var handshake Handshake
	if err := handshake.ReadFromConn(conn); err != nil {
		fmt.Println("Error reading handshake: ", err)
	}

	// handshake.PrettyPrint()

	switch handshake.NextState {
	case STATUS:
		if err := ReadStatusRequest(conn); err != nil {
			fmt.Println("Error reading status request: ", err)
		}

		statusResp := getStatus(conn)
		if err := statusResp.WriteToConn(conn); err != nil {
			fmt.Println("Error writing status response: ", err)
		}

		if err := PingPong(conn); err != nil {
			fmt.Println("Error with ping pong: ", err)
		}
	case LOGIN:
		fmt.Println("Login as attempted")

		var loginStart LoginStart
		if err := loginStart.ReadFromConn(conn); err != nil {
			fmt.Println("Error reading login start: ", err)
		}

		disconnect := getMsg(conn, loginStart)
		if err := disconnect.WriteToConn(conn); err != nil {
			fmt.Println("Error writing disconnect: ", err)
		}
	}
}
