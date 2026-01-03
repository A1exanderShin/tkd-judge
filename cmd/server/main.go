package main

import (
	"log"
	"net/http"
	httpapi "tkd-judge/internal/http"

	"tkd-judge/internal/ws"
)

func main() {
	hub := ws.NewHub()
	go hub.Run()

	http.Handle("/ws", ws.NewWSHandler(hub))
	http.Handle("/protocol", httpapi.NewProtocolHandler(hub))
	http.Handle("/protocol/pdf", httpapi.NewProtocolPDFHandler(hub))

	log.Println("server started on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
