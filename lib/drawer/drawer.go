package drawer

import (
	"bytes"
	"fmt"
	"math"
	"strings"

	svg "github.com/ajstarks/svgo"
	"github.com/lspaccatrosi16/go-cli-tools/logging"
	"github.com/lspaccatrosi16/go-libs/structures/cartesian"
	"github.com/lspaccatrosi16/tmappr/lib/types"
	"github.com/lspaccatrosi16/tmappr/lib/util"
)

const borderOffset = 100

func DrawMap(config *types.AppConfig, data *types.AppData) *bytes.Buffer {
	util.DebugSection("Drawing Map")

	logger := logging.GetLogger()

	buf := bytes.NewBuffer(nil)
	canvas := svg.New(buf)

	xR := float64(config.XRes) / float64(data.MaxX)
	yR := float64(config.YRes) / float64(data.MaxY)

	xc := scalefn(xR)
	yc := scalefn(yR)

	// stopRadius := int(1 * float64(config.Linewidth))
	halfWidth := int(0.5 * float64(config.Linewidth))
	threeeigthWidth := int(0.375 * float64(config.Linewidth))

	canvas.Start(config.XRes+borderOffset, config.YRes+borderOffset)

	logger.Debug("Drawing Line Segments")

	for _, segment := range data.Pathings.Segments {
		offset := func(n int) int {
			return config.Linewidth*n - (len(segment.Lines)-1)*halfWidth
		}

		segment.XOffsets = make([]int, len(segment.Lines))
		segment.YOffsets = make([]int, len(segment.Lines))

		var xo, yo int

		for i, line := range segment.Lines {
			if segment.Gradient.Coordinates()[1] != 0 {
				xo = offset(i)
			}

			if segment.Gradient.Coordinates()[0] != 0 {
				yo = offset(i)
			}

			segment.XOffsets[i] = xo
			segment.YOffsets[i] = yo
			lineStyle := fmt.Sprintf("stroke: %s; stroke-width: %dpx", line.Colour, config.Linewidth)
			circleStyle := fmt.Sprintf("fill: %s", line.Colour)

			xCoords := []int{xc(segment.Start[0]) + xo, xc(segment.End[0]) + xo}
			yCoords := []int{yc(segment.Start[1]) + yo, yc(segment.End[1]) + yo}
			canvas.Polyline(xCoords, yCoords, lineStyle)
			canvas.Circle(xc(segment.End[0])+xo, yc(segment.End[1])+yo, halfWidth, circleStyle)
		}
	}

	for c, stop := range data.CStopMap {
		var styleText string

		logger.Debug(fmt.Sprintf("Drawing Stop %s %s", stop.Name, c))
		lineSegment := data.Pathings.FindPrimarySegmentWithPoint(c)
		if lineSegment == nil {
			logger.Log(fmt.Sprintf("WARN: could not find line segment for stop %s", stop.Name))
			continue
		}

		segmentWidth := config.Linewidth * len(lineSegment.Lines)

		if stop.Type == types.AutoStopType {
			if len(data.StopLineMap[stop]) > 1 {
				stop.Type = types.Interchange
			} else {
				stop.Type = types.Normal
			}
		}

		switch stop.Type {
		case types.Interchange:
			var xs, ys int
			if len(lineSegment.Lines) > 1 {
				xs = segmentWidth
				ys = segmentWidth
			} else {
				xs = 3 * halfWidth
				ys = 3 * halfWidth
			}
			styleText = fmt.Sprintf("fill: #ffffff; stroke: #000000; stroke-width: %dpx", halfWidth)
			canvas.Roundrect(xc(c[0])-xs/2, yc(c[1])-ys/2, xs, ys, halfWidth, halfWidth, styleText)
		case types.Normal:
			for i := 0; i < len(stop.Lines); i++ {
				sl := data.LinesNames[stop.Lines[i]]
				for j, l := range lineSegment.Lines {
					if l.Code == sl.Code {
						xo := lineSegment.XOffsets[j]
						yo := lineSegment.YOffsets[j]
						canvas.Circle(xc(c[0])+xo, yc(c[1])+yo, threeeigthWidth, "fill: #ffffff")
						break
					}
				}
			}
		}

		textPlaced := false
		var selected cartesian.Direction

		if stop.Position != cartesian.NoDirection {
			selected = stop.Position
			textPlaced = true
		} else {
			for _, ca := range cartesian.CardinalPositions() {
				if checkAdjacent(data.UsedGrid, c, ca, lineSegment) {
					textPlaced = true
					selected = ca
					break
				}
			}
		}

		data.UsedGrid.Add(c.Add(selected.NextAcw().Coordinates()), 1)
		data.UsedGrid.Add(c.Add(selected.Coordinates()), 1)
		data.UsedGrid.Add(c.Add(selected.NextCw().Coordinates()), 1)

		if textPlaced {
			textStyleText := fmt.Sprintf("font-size: %.0fpx; text-anchor: middle;", 1.5*float64(config.Linewidth))
			words := strings.Split(stop.Name, " ")

			x := xc(c.Add(selected.Coordinates())[0])
			y := yc(c.Add(selected.Coordinates())[1]) + halfWidth - (len(words)/2)*(config.Linewidth+1)

			sx := selected.Coordinates()[0]
			sy := selected.Coordinates()[1]

			if sx > 0 {
				x += segmentWidth/2 + halfWidth
			} else if sx < 0 {
				x -= segmentWidth/2 + halfWidth
			}

			if sy > 0 {
				y += segmentWidth/2 + halfWidth
			} else if sy < 0 {
				y -= segmentWidth/2 + halfWidth
			}

			for i, word := range words {
				canvas.Text(x, y+i*2*(config.Linewidth+1), word, textStyleText)
			}

			textPlaced = true
		} else {
			logger.Log(fmt.Sprintf("WARN: could not place label %s", stop.Name))
		}
	}

	canvas.End()
	return buf
}

func scalefn(scale float64) func(c int) int {
	return func(c int) int {
		fv := float64(c)*math.Floor(scale) + 0.5*scale
		return int(math.Floor(fv)) + (borderOffset / 2)
	}
}

func checkAdjacent(grid *cartesian.CoordinateGrid[int], c cartesian.Coordinate, d cartesian.Direction, linesegment *types.CompoundSegment) bool {
	if linesegment.Gradient == d || linesegment.Gradient.Opposite() == d {
		return false
	}

	p1 := grid.Get(c.Add(d.NextAcw().Coordinates()))
	p2 := grid.Get(c.Add(d.Coordinates()))
	p3 := grid.Get(c.Add(d.NextCw().Coordinates()))

	return p1 == 0 && p2 == 0 && p3 == 0
}
