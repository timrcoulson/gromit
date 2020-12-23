package guitar

import (
	"fmt"
	"github.com/timrcoulson/gromit/guitar/arpeggio"
	"github.com/timrcoulson/gromit/guitar/scales"
	"math/rand"
)

type Guitar struct {

}

func (g *Guitar) Output() string {
	mode := scales.Mode(rand.Intn(7))
	arpeggioType := arpeggio.ArpeggioPosition(rand.Intn(9))
	randomScale := scales.Scale(rand.Intn(12), mode)
	randomArpeggio := arpeggio.Arpeggio(rand.Intn(12), arpeggioType)

	output := fmt.Sprintf("# Scale of the day - %s \n%s", mode.String(), randomScale)
	output += fmt.Sprintf("# Arpeggio of the day - %s \n%s", arpeggioType.String(), randomArpeggio)
	output += "\n# Arpeggio of the day\n" + randomArpeggio

	return output
}
