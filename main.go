package main

import (
	"fmt"
	"net"

	"github.com/nathanielfernandes/motd/slp"
)

func onStatus(conn net.Conn) (slp.StatusResponse, error) {
	fmt.Println("status")
	return slp.GrabStatus("mc.nathanferns.xyz:25565")
}

func onLogin(conn net.Conn, loginStart slp.LoginStart) error {
	fmt.Println("login")

	return nil
}

func onError(conn net.Conn, err error) {
	fmt.Println(err)
}

func main() {
	server := slp.NewServer().
		OnStatus(onStatus).
		OnLogin(onLogin).
		OnError(onError)

	fmt.Println("Server started")
	server.ListenAndServe(":25565")
}
