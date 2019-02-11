package main

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/websocket"
)

type client struct {
	name       string
	clientConn *websocket.Conn
}

// Clients Connected
var clients []client

func echoHandler(w http.ResponseWriter, r *http.Request) {
	conn, _ := upgrader.Upgrade(w, r, nil)
	for {
		msgType, msg, err := conn.ReadMessage()
		if err != nil {
			return
		}
		for i := 0; i < len(clients); i++ {
			if conn.RemoteAddr() == clients[i].clientConn.RemoteAddr() {
				fmt.Printf("%s sent: %s\n", clients[i].name, string(msg))
			}
		}
		if string(msg) == "SOCKET_OPEN" {
			c := client{"DEFAULT", conn}
			clients = append(clients, c)
			fmt.Println("client has joined")
			fmt.Println("clients len:" + strconv.Itoa(len(clients)))
		} else if string(msg) == "SOCKET_CLOSED" {
			for i := 0; i < len(clients); i++ {
				if clients[i].clientConn.RemoteAddr() == conn.RemoteAddr() {
					fmt.Println("Removing " + clients[i].name + " from slice")
					clients = append(clients[:i], clients[i+1:]...)
				}
			}
			fmt.Println("clients len:" + strconv.Itoa(len(clients)))
		} else if strings.Contains(string(msg), "SET NAME:") {
			for i := 0; i < len(clients); i++ {
				if clients[i].clientConn.RemoteAddr() == conn.RemoteAddr() {
					clients[i].name = strings.TrimPrefix(string(msg), "SET NAME:")
				}
			}
		} else if err = conn.WriteMessage(msgType, msg); err != nil {
			return
		} else {
			clientName := ""
			fmt.Print(clientName)
			for i := 0; i < len(clients); i++ {
				if clients[i].clientConn.RemoteAddr() == conn.RemoteAddr() {
					clientName = clients[i].name
				}
			}
			for i := 0; i < len(clients); i++ {
				if clients[i].clientConn.RemoteAddr() != conn.RemoteAddr() {
					fullMsg := clientName + ": " + string(msg)
					clients[i].clientConn.WriteMessage(msgType, []byte(fullMsg))
				}
			}
		}
	}
}

func chatServerHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "/home/haxxionlaptop/Documents/VisualStudioCode/Go/WebWiki/html/chatServer.html")
}
