package main

import (
	"log"
	"net/http"

	"tkd-judge/internal/ws"
)

func main() {
	hub := ws.NewHub()

	go hub.Run()

	wsHandler := ws.NewWSHandler(hub)

	http.Handle("/ws", wsHandler)

	log.Println("server started on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
