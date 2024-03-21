package drawer

import (
	"bytes"
	"fmt"
	"math"

	svg "github.com/ajstarks/svgo"
	"github.com/lspaccatrosi16/tmappr/lib/types"
)

const borderOffset = 20

func DrawMap(config *types.AppConfig, pathings []*types.PathedLine, maxX, maxY int) *bytes.Buffer {
	buf := bytes.NewBuffer(nil)
	canvas := svg.New(buf)

	xR := float64(config.XRes) / float64(maxX)
	yR := float64(config.YRes) / float64(maxY)

	xc := scalefn(xR)
	yc := scalefn(yR)

	canvas.Start(config.XRes+borderOffset, config.YRes+borderOffset)

	for _, path := range pathings {
		lineStyle := fmt.Sprintf("stroke: %s; stroke-width: %dpx", path.Line.Colour, config.Linewidth)
		for _, segment := range path.Segments {
			xCoords := []int{xc(segment.Start[0]), xc(segment.End[0])}
			yCoords := []int{yc(segment.Start[1]), yc(segment.End[1])}

			canvas.Polyline(xCoords, yCoords, lineStyle)
		}
	}

	canvas.End()
	return buf
}

func scalefn(scale float64) func(c int) int {
	return func(c int) int {
		fv := float64(c) * math.Floor(scale)
		return int(math.Floor(fv)) + (borderOffset / 2)
	}
}
