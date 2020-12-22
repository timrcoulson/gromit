package guitar

import (
	"github.com/timrcoulson/gromit/guitar/arpeggio"
	"github.com/timrcoulson/gromit/guitar/scales"
	"math/rand"
)

type Guitar struct {

}

func (g *Guitar) Output() string {
	randomScale := scales.Scale(rand.Intn(12), scales.Mode(rand.Intn(7)))
	randomArpeggio := arpeggio.Arpeggio(rand.Intn(12), arpeggio.ArpeggioPosition(rand.Intn(4)))

	output := "# Scale of the day\n" + randomScale
	output += "\n# Arpeggio of the day\n" + randomArpeggio

	return output
}
