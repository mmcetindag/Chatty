package server

import (
	"../chatty"
	"encoding/gob"

	"fmt"
	"net"
)

var users = make(map[string]chatty.ChatConn)

func handleConnection(conn net.Conn) {
	userConn := chatty.ChatConn{
		Enc:  gob.NewEncoder(conn),
		Dec:  gob.NewDecoder(conn),
		Conn: conn,
	}

	initialMsg, _ := chatty.RecvMsg(userConn)

	if _, ok := users[initialMsg.Username]; ok {
		fmt.Printf("User from '%s' tried to connect with duplicate name '%s'\n", userConn.Conn.RemoteAddr().String(), initialMsg.Username)
		msgSend := chatty.ChattyMsg{Body: "Duplicate name!", Action: chatty.ERROR}
		chatty.SendMsg(userConn, msgSend)
		conn.Close()
		return
	}

	fmt.Printf("User %s succesfully connected with IP %s\n", initialMsg.Username, userConn.Conn.RemoteAddr().String())
	users[initialMsg.Username] = userConn

	for {
		msg, _ := chatty.RecvMsg(userConn)

		switch msg.Action {
		case chatty.DISCONNECT:
			fmt.Printf("User %s disconnected\n", initialMsg.Username)
			delete(users, initialMsg.Username)
			conn.Close()
			return
		case chatty.LIST:
			usersText := "List of users:\n"
			for user := range users {
				usersText += fmt.Sprintf(" - %s\n", user)
			}
			msgSend := chatty.ChattyMsg{Body: usersText, Action: chatty.LIST}
			chatty.SendMsg(userConn, msgSend)
		case chatty.MSG:
			if _, ok := users[msg.Username]; !ok {
				msgSend := chatty.ChattyMsg{Body: "Recipient doesn't exist", Action: chatty.ERROR}
				chatty.SendMsg(userConn, msgSend)
			} else {
				msgSend := chatty.ChattyMsg{Username: initialMsg.Username, Body: msg.Body, Action: chatty.MSG}
				chatty.SendMsg(users[msg.Username], msgSend)
			}
		}
	}
}

func Start() {
	listen, port, err := chatty.OpenListener()
	fmt.Printf("Listening on port %v\n", port)

	if err != nil {
		fmt.Println(err)
		return
	}

	for {
		conn, err := listen.Accept() // this blocks until connection or error
		if err != nil {
			fmt.Println(err)
			continue
		}
		go handleConnection(conn)
	}
}
