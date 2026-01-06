package main

import (
	"log"
	"net/http"

	"tkd-judge/internal/discipline"
	"tkd-judge/internal/ws"
)

func main() {
	// 1. Выбор дисциплины (composition root)
	// Сейчас используем бой
	d := discipline.NewFightDiscipline()

	// 2. Центральный хаб
	hub := ws.NewHub(d)
	go hub.Run()

	// 3. WebSocket endpoint
	http.Handle("/ws", ws.NewWSHandler(hub))

	// 4. UI (статические файлы)
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
