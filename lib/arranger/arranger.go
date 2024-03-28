package arranger

import (
	"fmt"
	"strings"

	"github.com/lspaccatrosi16/tmappr/lib/types"
)

type genLine struct {
	Line     *types.Line
	Segments []*types.CompoundSegment
	WSegs    []weightSeg
}

type weightSeg struct {
	Start   int
	End     int
	Segment *types.CompoundSegment
}

func (g *genLine) String() string {
	segText := []string{}
	if len(g.WSegs) > 0 {
		for _, s := range g.WSegs {
			segText = append(segText, s.Segment.LineSegment.String()+fmt.Sprintf(" %d:%d", s.Start, s.End))
		}
	} else {
		for _, s := range g.Segments {
			segText = append(segText, s.LineSegment.String())
		}
	}
	return fmt.Sprintf("%s: %s", g.Line.Code, strings.Join(segText, ", "))
}

func Arrange(config *types.AppConfig, data *types.AppData) {
	idivlines := regenLines(data)
	calculateWeights(idivlines)

	for _, line := range idivlines {
		fmt.Println(line.String())
	}
}

func calculateWeights(lines []*genLine) {
	for _, line := range lines {
		wsegs := []weightSeg{}
		for i := 0; i < len(line.Segments); i++ {
			seg := line.Segments[i]
			wseg := weightSeg{
				Segment: seg,
			}
			if i > 0 {
				var weight int
				weighta := seg.Gradient.NumberCw(line.Segments[i-1].Gradient)
				weightb := seg.Gradient.NumberAcw(line.Segments[i-1].Gradient)
				if weighta < weightb {
					weight = weighta
				} else {
					weight = weightb
				}

				if weight > 4 {
					weight = 4 - weight
				}
				wseg.Start = weight
			}

			if i < len(line.Segments)-1 {
				var weight int
				weighta := seg.Gradient.NumberCw(line.Segments[i+1].Gradient)
				weightb := seg.Gradient.NumberAcw(line.Segments[i+1].Gradient)
				if weighta < weightb {
					weight = weighta
				} else {
					weight = weightb
				}

				if weight > 4 {
					weight = 4 - weight
				}
				wseg.End = weight
			}

			wsegs = append(wsegs, wseg)
		}
		line.WSegs = wsegs
	}
}

func regenLines(data *types.AppData) []*genLine {
	lines := []*genLine{}

	for _, line := range data.Lines {
		gl := genLine{
			Line:     line,
			Segments: []*types.CompoundSegment{},
		}
		curseg := data.Pathings.FindSegmentWithPointLine(line.Stops[0].IntCoordinates, line)[0]

		i := 0
		for {
			if len(gl.Segments) == 0 || gl.Segments[len(gl.Segments)-1] != curseg {
				gl.Segments = append(gl.Segments, curseg)
			}
			i++
			if i >= len(line.Stops) {
				break
			}

			if curseg.Start == line.Stops[i].IntCoordinates {
				newSegs := data.Pathings.FindSegmentWithPointLine(line.Stops[i].IntCoordinates, line)
				for _, s := range newSegs {
					if s.End != curseg.End {
						curseg = s
						break
					}
				}
			} else {
				newSegs := data.Pathings.FindSegmentWithPointLine(line.Stops[i].IntCoordinates, line)
				for _, s := range newSegs {
					if s.Start != curseg.Start {
						curseg = s
						break
					}
				}
			}
		}
		lines = append(lines, &gl)
	}

	return lines

}
