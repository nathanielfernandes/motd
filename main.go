package main

import (
	"fmt"
	"net"

	"github.com/nathanielfernandes/motd/slp"
)

func main() {
	ln, err := net.Listen("tcp", ":25565")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Listening on :25565")
	defer ln.Close()

	favicon, _ := slp.FaviconFromFile("favicon.png")
	genStatus := func(conn net.Conn) slp.StatusResponse {
		return slp.StatusResponse{
			Version: slp.Version{
				Name:     "1.19.3",
				Protocol: 761,
			},
			// or
			// Version: slp.SpoofVersion("blah blah blah")
			Players: slp.Players{
				Max:    100,
				Online: 7,
			},
			Description: slp.Description{
				Text: "\u00a7bHello \u00a7aWorld \u00a7r\u00a7ka \u00a7rfrom: " + conn.RemoteAddr().String(),
			},
			Favicon:            favicon,
			PreviewsChat:       true,
			EnforcesSecureChat: true,
		}
	}

	genDisconnect := func(conn net.Conn, loginStart slp.LoginStart) slp.Disconnect {
		msg := slp.Chat{
			Text:  "Hey " + loginStart.Name + ",\n",
			Color: "white",
			Extra: []slp.Chat{
				{
					Text:  "this is not a real server\n",
					Color: "dark_gray",
				},
				{
					Text:  "please go away",
					Color: "red",
					Font:  "minecraft:alt",
				},
			},
		}
		return slp.DisconnectWithChatMsg(msg)
	}

	for {
		conn, err := ln.Accept()

		if err != nil {
			fmt.Println(err)
		}

		go slp.HandleConnection(conn, genStatus, genDisconnect)
	}
}
