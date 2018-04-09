package client

import (
	"../chatty"

	"fmt"
	"bufio"
	"os"
	"strings"
)

func receiveMessages(conn chatty.ChatConn) {
	for {
		msg, _ := chatty.RecvMsg(conn)

		switch msg.Action {
		case chatty.ERROR:
			fmt.Printf("ERROR: %s\n", msg.Body)
		case chatty.LIST:
			fmt.Print(msg.Body)
		case chatty.MSG:
			fmt.Printf("%s: %s", msg.Username, msg.Body)
		}
	}
}

func Start(user string, serverPort string, serverAddr string) {
	// Connect to chat server
	chatConn, err := chatty.ServerConnect(user, serverAddr, serverPort)
	if err != nil {
		fmt.Printf("unable to connect to server: %v\n", err)
		return
	}

	fmt.Printf("Connected to server at %v\n"+
		"Type '/list' to get a list of users\n"+
		"Type '/disconnect' to... well... disconnect\n"+
		"Type a username to send that user a message\n", serverAddr)

	go receiveMessages(chatConn) // Use another thread to receive messages

	reader := bufio.NewReader(os.Stdin)
	for {
		text, _ := reader.ReadString('\n')
		text = strings.TrimSpace(text)

		var msg chatty.ChattyMsg
		switch text {
		case "/list":
			msg = chatty.ChattyMsg{Action: chatty.LIST}
		case "/disconnect":
			msg = chatty.ChattyMsg{Action: chatty.DISCONNECT}
		default:
			fmt.Printf("Write the message you want to send to %s\n", text)
			body, _ := reader.ReadString('\n')
			msg = chatty.ChattyMsg{Username: text, Body: body, Action: chatty.MSG}
		}

		chatty.SendMsg(chatConn, msg)

		if msg.Action == chatty.DISCONNECT {
			os.Exit(0)
		}
	}
}
