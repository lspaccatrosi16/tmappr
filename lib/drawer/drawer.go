package drawer

import (
	"bytes"
	"fmt"
	"math"

	svg "github.com/ajstarks/svgo"
	"github.com/lspaccatrosi16/go-libs/structures/cartesian"
	"github.com/lspaccatrosi16/tmappr/lib/types"
)

const borderOffset = 100

func DrawMap(config *types.AppConfig, pathings *map[string]*types.PathedLine, cStopMap *map[cartesian.Coordinate]*types.Stop, maxX, maxY int) *bytes.Buffer {
	buf := bytes.NewBuffer(nil)
	canvas := svg.New(buf)

	xR := float64(config.XRes) / float64(maxX)
	yR := float64(config.YRes) / float64(maxY)

	xc := scalefn(xR)
	yc := scalefn(yR)

	stopRadius := int(1 * float64(config.Linewidth))
	halfWidth := int(0.5 * float64(config.Linewidth))
	quarterWidth := int(0.25 * float64(config.Linewidth))

	canvas.Start(config.XRes+borderOffset, config.YRes+borderOffset)

	for _, path := range *pathings {
		lineStyle := fmt.Sprintf("stroke: %s; stroke-width: %dpx", path.Line.Colour, config.Linewidth)
		circleStyle := fmt.Sprintf("fill: %s", path.Line.Colour)
		for _, segment := range path.Segments {
			xCoords := []int{xc(segment.Start[0]), xc(segment.End[0])}
			yCoords := []int{yc(segment.Start[1]), yc(segment.End[1])}

			canvas.Polyline(xCoords, yCoords, lineStyle)
			canvas.Circle(xc(segment.End[0]), yc(segment.End[1]), halfWidth, circleStyle)
		}
	}
	for c, stop := range *cStopMap {
		var styleText string
		if len(stop.Lines) > 1 {
			styleText = fmt.Sprintf("fill: #ffffff; stroke: #000000; stroke-width: %dpx", halfWidth)
			canvas.Circle(xc(c[0]), yc(c[1]), stopRadius, styleText)
		} else if len(stop.Lines) == 1 {
			pathing := (*pathings)[stop.Lines[0]]
			styleText = fmt.Sprintf("fill: #ffffff; stroke-width: %dpx; stroke: %s;", stopRadius, pathing.Line.Colour)
			canvas.Circle(xc(c[0]), yc(c[1]), quarterWidth, styleText)
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
