package httpapi

import (
	"net/http"

	"tkd-judge/internal/pdf"
	"tkd-judge/internal/ws"
)

type ProtocolPDFHandler struct {
	hub *ws.Hub
}

func NewProtocolPDFHandler(hub *ws.Hub) *ProtocolPDFHandler {
	return &ProtocolPDFHandler{hub: hub}
}

func (h *ProtocolPDFHandler) ServeHTTP(w http.ResponseWriter, _ *http.Request) {
	proto := h.hub.BuildProtocol()

	data, err := pdf.BuildProtocolPDF(proto)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/pdf")
	w.Header().Set("Content-Disposition", "attachment; filename=protocol.pdf")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(data)
}
