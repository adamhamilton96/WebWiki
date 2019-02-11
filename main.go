package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/gorilla/websocket"
)

// WebSocket
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	files, _ := ioutil.ReadDir("/home/haxxionlaptop/Documents/VisualStudioCode/Go/WebWiki/pages/")
	var txtFiles []string
	for _, f := range files {
		if strings.Contains(f.Name(), ".txt") {
			fileName := f.Name()
			name := strings.TrimSuffix(fileName, filepath.Ext(fileName))
			txtFiles = append(txtFiles, name)
		}
	}
	templates.ExecuteTemplate(w, "home.html", txtFiles)
}

func renderTemplate(w http.ResponseWriter, tmpl string, p *Page) {
	err := templates.ExecuteTemplate(w, tmpl+".html", p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func makeHandler(fn func(http.ResponseWriter, *http.Request, string)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		m := validPath.FindStringSubmatch(r.URL.Path)
		if m == nil {
			http.NotFound(w, r)
			return
		}
		fn(w, r, m[2])
	}
}

func main() {
	fmt.Println("Wiki Web Server Online")

	// Wiki
	http.HandleFunc("/view/", makeHandler(viewHandler))
	http.HandleFunc("/edit/", makeHandler(editHandler))
	http.HandleFunc("/save/", makeHandler(saveHandler))

	// Home
	http.HandleFunc("/", rootHandler)

	// Chat Server
	http.HandleFunc("/chatserver/", chatServerHandler)
	http.HandleFunc("/echo", echoHandler)

	// Shared Canvas
	http.HandleFunc("/sharedcanvas/", sharedCanvasHandler)
	http.HandleFunc("/canvasecho", canvasEchoHandler)

	// Snakes and Ladders
	http.HandleFunc("/snakes+ladders/", snakesAndLaddersHandler)
	http.HandleFunc("/snakes+laddersecho", snakesAndLaddersEchoHandler)

	log.Fatal(http.ListenAndServe("192.168.42.2:1255", nil))
}
