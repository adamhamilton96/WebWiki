package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/websocket"
)

// Painters Connected
var painters []*websocket.Conn

func canvasEchoHandler(w http.ResponseWriter, r *http.Request) {
	conn, _ := upgrader.Upgrade(w, r, nil)
	for {
		msgType, msg, err := conn.ReadMessage()
		if err != nil {
			return
		}
		if string(msg) == "SOCKET_OPEN" {
			fmt.Println("painter has joined")
			painters = append(painters, conn)
			fmt.Println("painters len:" + strconv.Itoa(len(painters)))
		} else if string(msg) == "SOCKET_CLOSED" {
			for i := 0; i < len(painters); i++ {
				if painters[i].RemoteAddr() == conn.RemoteAddr() {
					fmt.Printf("%s removed from slice\n", painters[i].RemoteAddr())
					painters = append(painters[:i], painters[i+1:]...)
				}
			}
			fmt.Println("painters len:" + strconv.Itoa(len(painters)))
		} else if err = conn.WriteMessage(msgType, msg); err != nil {
			return
		} else {
			for i := 0; i < len(painters); i++ {
				if painters[i].RemoteAddr() != conn.RemoteAddr() {
					painters[i].WriteMessage(msgType, msg)
				}
			}
		}
	}
}

func sharedCanvasHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "/home/haxxionlaptop/Documents/VisualStudioCode/Go/WebWiki/html/sharedCanvas.html")
}
