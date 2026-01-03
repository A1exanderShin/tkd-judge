package httpapi

import (
	"encoding/json"
	"net/http"

	"tkd-judge/internal/ws"
)

type ProtocolHandler struct {
	hub *ws.Hub
}

func NewProtocolHandler(hub *ws.Hub) *ProtocolHandler {
	return &ProtocolHandler{hub: hub}
}

func (h *ProtocolHandler) ServeHTTP(w http.ResponseWriter, _ *http.Request) {
	proto := h.hub.BuildProtocol()

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(proto)
}
