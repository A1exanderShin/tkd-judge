package main

import (
	"log"
	"net/http"

	"tkd-judge/internal/discipline"
	"tkd-judge/internal/ws"
)

func main() {

	// ===== DISCIPLINES =====

	fight := discipline.NewFightDiscipline()

	pattern := discipline.NewPatternDiscipline(
		[]string{
			"technique",
			"balance",
			"power",
			"rhythm",
			"expression",
		},
		5, // количество судей
	)

	// ===== ROUTER =====

	router := discipline.NewRouter(fight)

	// ===== HUB =====

	hub := ws.NewHub(
		router,
		fight,
		pattern,
	)

	go hub.Run()

	// ===== WEBSOCKET =====

	http.Handle("/ws", ws.NewWSHandler(hub))

	// ===== UI =====

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
