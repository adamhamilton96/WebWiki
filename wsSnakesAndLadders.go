package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/websocket"
)

type snadder struct {
	headx int
	heady int
	tailx int
	taily int
}

type player struct {
	x          int
	y          int
	turn       bool
	right      bool
	clientConn *websocket.Conn
}

var players []player
var board [][]int
var ladders []snadder
var snakes []snadder
var gridSize = 10
var gameReady = false

func snakesAndLaddersEchoHandler(w http.ResponseWriter, r *http.Request) {
	conn, _ := upgrader.Upgrade(w, r, nil)
	for {
		msgType, msg, err := conn.ReadMessage()
		if err != nil {
			return
		}
		if string(msg) == "SOCKET_OPEN" {
			if len(players) < 2 {
				newPlayer := player{0, 9, false, true, conn}
				players = append(players, newPlayer)
				fmt.Println("players has joined")
				fmt.Println("players len:" + strconv.Itoa(len(players)))
			}
			if len(players) == 2 {
				gameReady = true
				setup()
				for i := 0; i < len(players); i++ {
					players[i].clientConn.WriteMessage(msgType, []byte(setupMsg()))
				}
			} else {
				gameReady = false
			}
		} else if string(msg) == "SOCKET_CLOSED" {
			for i := 0; i < len(players); i++ {
				if players[i].clientConn.RemoteAddr() == conn.RemoteAddr() {
					fmt.Printf("%s removed from slice\n", players[i].clientConn.RemoteAddr())
					players = append(players[:i], players[i+1:]...)
					gameReady = false
					fmt.Println("players len:" + strconv.Itoa(len(players)))
				}
			}
		} else if err = conn.WriteMessage(msgType, msg); err != nil {
			return
		} else {
			for i := 0; i < len(players); i++ {
				if players[i].clientConn.RemoteAddr() != conn.RemoteAddr() {
					players[i].clientConn.WriteMessage(msgType, msg)
				}
			}
		}
	}
}

func snakesAndLaddersHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "/home/haxxionlaptop/Documents/VisualStudioCode/Go/WebWiki/html/snakesAndLadders.html")
}

func setup() {
	// Set player positions
	for i := 0; i < len(players); i++ {
		players[i].x = 0
		players[i].y = 9
		players[i].turn = false
		players[i].right = true
	}
	players[0].turn = true

	// Make board
	board = make([][]int, gridSize)
	for i := range board {
		board[i] = make([]int, gridSize)
	}

	// Snakes
	snake1 := newSnadder(1, 0, 2, 2)
	snake2 := newSnadder(5, 0, 5, 2)
	snake3 := newSnadder(6, 1, 3, 7)
	snake4 := newSnadder(3, 3, 0, 5)
	snakes = append(snakes, snake1)
	snakes = append(snakes, snake2)
	snakes = append(snakes, snake3)
	snakes = append(snakes, snake4)

	// Ladders
	ladder1 := newSnadder(9, 0, 9, 2)
	ladder2 := newSnadder(3, 1, 7, 7)
	ladder3 := newSnadder(0, 1, 2, 3)
	ladder4 := newSnadder(9, 6, 8, 9)
	ladders = append(ladders, ladder1)
	ladders = append(ladders, ladder2)
	ladders = append(ladders, ladder3)
	ladders = append(ladders, ladder4)
}

func setupMsg() string {
	setupMsg := " "
	setupMsg = "SETUP " + strconv.FormatBool(gameReady) + " " + strconv.Itoa(gridSize) + " " + strconv.Itoa(len(players)) + " "
	for i := 0; i < len(players); i++ {
		setupMsg += strconv.Itoa(players[i].x) + " " + strconv.Itoa(players[i].y) + " "
	}
	setupMsg += strconv.Itoa(len(snakes)) + " "
	for i := 0; i < len(snakes); i++ {
		setupMsg += strconv.Itoa(snakes[i].headx) + " " + strconv.Itoa(snakes[i].heady) + " " +
			strconv.Itoa(snakes[i].tailx) + " " + strconv.Itoa(snakes[i].taily) + " "
	}
	setupMsg += strconv.Itoa(len(ladders)) + " "
	for i := 0; i < len(ladders); i++ {
		setupMsg += strconv.Itoa(ladders[i].headx) + " " + strconv.Itoa(ladders[i].heady) + " " +
			strconv.Itoa(ladders[i].tailx) + " " + strconv.Itoa(ladders[i].taily) + " "
	}
	return setupMsg
}

// Snake and Ladder Constructor
func newSnadder(iheadx int, iheady int, itailx int, itaily int) snadder {
	return snadder{iheadx, iheady, itailx, itaily}
}

// Dice Rolling and Movement Function
func roll() {
	if gameReady == true {
		rand.Seed(time.Now().UnixNano())
		roll := (rand.Intn(6) + 1)

		// Find players turn
		for i := 0; i < len(players); i++ {
			if players[i].turn == true {
				for roll > 0 {
					if players[i].right == true && players[i].x < 9 && roll < 9-players[i].x {
						players[i].x += roll
						roll = 0
					} else if players[i].right == false && players[i].x > 0 && roll < 0+players[i].x {
						players[i].x -= roll
						roll = 0
					} else if players[i].right == true && players[i].x < 9 {
						for roll > 0 && players[i].x != 9 {
							players[i].x++
							roll--
						}
					} else if players[i].right == false && players[i].x > 0 {
						if players[i].y == 0 && roll-players[i].x > 0 {
							roll = 0
						} else {
							for roll > 0 && players[i].x > 0 {
								players[i].x--
								roll--
							}
						}
					} else {
						players[i].y--
						roll++
						if players[i].y%2 == 0 {
							players[i].x = 9
							players[i].right = false
						} else {
							players[i].x = 0
							players[i].right = true
						}
					}
					for j := 0; j < len(snakes); j++ {
						snake := snakes[j]
						if players[i].x == snake.headx && players[i].y == snake.heady {
							players[i].x = snake.tailx
							players[i].y = snake.taily
							roll = 0
							if players[i].y%2 == 0 {
								players[i].right = false
							} else {
								players[i].right = true
							}
						}
					}
					for j := 0; j < len(ladders); j++ {
						ladder := ladders[j]
						if players[i].x == ladder.tailx && players[i].y == ladder.taily {
							players[i].x = ladder.headx
							players[i].y = ladder.heady
							roll = 0
							if players[i].y%2 == 0 {
								players[i].right = false
							} else {
								players[i].right = true
							}
						}
					}
					roll--
					fmt.Println(strconv.Itoa(roll) + " p[" + strconv.Itoa(i) + "].x: " + strconv.Itoa(players[i].x) +
						" p[" + strconv.Itoa(i) + "].y: " + strconv.Itoa(players[i].y))
				}
				if players[i].x == 0 && players[i].y == 0 {
					fmt.Println("Player " + strconv.Itoa(i) + " Winner")
					// TODO GAME ENDED
				}
				// Next players turn
				players[i].turn = false
				if i+1 > len(players) {
					players[0].turn = true
				} else {
					players[i+1].turn = true
				}
				return
			}
		}
	}
}
