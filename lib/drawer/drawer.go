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

func DrawMap(config *types.AppConfig, pathings *types.PathedSystem, cStopMap *map[cartesian.Coordinate]*types.Stop, lineMap *map[string]*types.Line, combinedGrid *cartesian.CoordinateGrid[bool], maxX, maxY int) *bytes.Buffer {
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

	for _, segment := range pathings.Segments {
		offset := func(n int) int {
			return config.Linewidth*n - (len(segment.Lines)-1)*halfWidth
		}

		for i, line := range segment.Lines {
			o := offset(i)
			lineStyle := fmt.Sprintf("stroke: %s; stroke-width: %dpx", line.Colour, config.Linewidth)
			circleStyle := fmt.Sprintf("fill: %s", line.Colour)

			xCoords := []int{xc(segment.Start[0]) + o, xc(segment.End[0]) + o}
			yCoords := []int{yc(segment.Start[1]) + o, yc(segment.End[1]) + o}
			canvas.Polyline(xCoords, yCoords, lineStyle)
			canvas.Circle(xc(segment.End[0])+o, yc(segment.End[1])+o, halfWidth, circleStyle)
		}
	}
	for c, stop := range *cStopMap {
		var styleText string
		if len(stop.Lines) > 1 {
			styleText = fmt.Sprintf("fill: #ffffff; stroke: #000000; stroke-width: %dpx", halfWidth)
			canvas.Circle(xc(c[0]), yc(c[1]), stopRadius, styleText)
		} else if len(stop.Lines) == 1 {
			line := (*lineMap)[stop.Lines[0]]
			styleText = fmt.Sprintf("fill: #ffffff; stroke-width: %dpx; stroke: %s;", stopRadius, line.Colour)
			canvas.Circle(xc(c[0]), yc(c[1]), quarterWidth, styleText)
		}

		combos := []cartesian.Coordinate{
			{1, 0},
			{1, 1},
			{1, -1},
			{0, 1},
			{0, -1},
			{-1, 0},
			{-1, 1},
			{-1, -1},
		}

		for _, ca := range combos {
			if !combinedGrid.Get(ca) {
				newCoordinate := c.Add(ca)
				combinedGrid.Add(newCoordinate, true)

				textStyleText := fmt.Sprintf("font-size: %dpx", 2*config.Linewidth)

				canvas.Text(xc(newCoordinate[0]), yc(newCoordinate[0]), stop.Name, textStyleText)
				break
			}
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
