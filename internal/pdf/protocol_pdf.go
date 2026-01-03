package pdf

import (
	"bytes"
	"fmt"
	"time"

	"tkd-judge/internal/protocol"

	"github.com/jung-kurt/gofpdf"
)

func BuildProtocolPDF(p protocol.Protocol) ([]byte, error) {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.SetTitle("TKD Fight Protocol", false)
	pdf.AddPage()

	pdf.SetFont("Arial", "B", 16)
	pdf.Cell(0, 10, "TKD Fight Protocol")
	pdf.Ln(12)

	pdf.SetFont("Arial", "", 11)
	pdf.Cell(0, 8, fmt.Sprintf("Generated: %s", p.GeneratedAt.Format(time.RFC3339)))
	pdf.Ln(6)
	pdf.Cell(0, 8, fmt.Sprintf("State: %s", p.State))
	pdf.Ln(10)

	pdf.SetFont("Arial", "B", 12)
	pdf.Cell(0, 8, "Score")
	pdf.Ln(8)
	pdf.SetFont("Arial", "", 11)
	pdf.Cell(0, 8, fmt.Sprintf("Red: %d   Blue: %d", p.Score.Red, p.Score.Blue))
	pdf.Ln(10)

	pdf.SetFont("Arial", "B", 12)
	pdf.Cell(0, 8, "Warnings")
	pdf.Ln(8)
	pdf.SetFont("Arial", "", 11)
	pdf.Cell(0, 8, fmt.Sprintf("Red: %d   Blue: %d", p.Warnings.Red, p.Warnings.Blue))
	pdf.Ln(10)

	pdf.SetFont("Arial", "B", 12)
	pdf.Cell(0, 8, "Events")
	pdf.Ln(8)

	pdf.SetFont("Arial", "", 10)
	if len(p.Events) == 0 {
		pdf.Cell(0, 8, "No events")
		pdf.Ln(8)
	} else {
		for _, e := range p.Events {
			pdf.MultiCell(0, 6, fmt.Sprintf("%v", e), "", "", false)
			pdf.Ln(1)
		}
	}

	var buf bytes.Buffer
	if err := pdf.Output(&buf); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
