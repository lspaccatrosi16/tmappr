package arranger

import (
	"fmt"
	"strings"

	"github.com/lspaccatrosi16/tmappr/lib/types"
)

type genLine struct {
	Line     *types.Line
	Segments []*types.CompoundSegment
}

func (g *genLine) String() string {
	segText := []string{}
	for _, s := range g.Segments {
		segText = append(segText, s.LineSegment.String())
	}
	return fmt.Sprintf("%s: %s", g.Line.Code, strings.Join(segText, ", "))
}

func Arrange(config *types.AppConfig, data *types.AppData) {
	idivlines := regenLines(data)

	for _, line := range idivlines {
		fmt.Println(line.String())
	}
}

func regenLines(data *types.AppData) []genLine {
	lines := []genLine{}

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
		lines = append(lines, gl)
	}

	return lines

}
