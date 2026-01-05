package ws

import (
	"net/http"
	"strconv"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

func NewWSHandler(hub *Hub) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// --- role ---
		roleStr := r.URL.Query().Get("role")

		var role Role
		switch roleStr {
		case "main":
			role = RoleMainJudge
		default:
			role = RoleSideJudge
		}

		// --- upgrade ---
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			return
		}

		client := NewClient(conn, role)

		// --- назначение judgeID (ТОЛЬКО для боковых судей) ---
		if role == RoleSideJudge {
			judgeParam := r.URL.Query().Get("judge")

			var id int
			if judgeParam != "" {
				id, err = strconv.Atoi(judgeParam)
				if err != nil || id < 1 || id > hub.cfg.JudgesCount {
					conn.Close()
					return
				}

				// если судья уже занят — не пускаем
				if _, exists := hub.sideJudges[id]; exists {
					conn.Close()
					return
				}
			} else {
				id = hub.nextFreeJudgeID()
				if id == 0 {
					conn.Close()
					return
				}
			}

			client.judgeID = id
		}

		// --- регистрация клиента ---
		hub.register <- client

		// --- read loop ---
		go func() {
			defer func() {
				hub.unregister <- client
				client.close()
			}()

			for {
				var msg Event
				if err := conn.ReadJSON(&msg); err != nil {
					return
				}
				hub.Publish(msg, client)
			}
		}()
	}
}
