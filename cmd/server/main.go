package main

import (
	"log"
	"net/http"

	"tkd-judge/internal/ws"
)

func main() {
	// центральный хаб
	hub := ws.NewHub()
	go hub.Run()

	// WebSocket
	http.Handle("/ws", ws.NewWSHandler(hub))

	// UI (ВАЖНО)
	http.Handle(
		"/ui/",
		http.StripPrefix(
			"/ui/",
			http.FileServer(http.Dir("ui")),
		),
	)

	log.Println("server started on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
