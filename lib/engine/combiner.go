package engine

import (
	"fmt"

	"github.com/lspaccatrosi16/go-cli-tools/logging"
	"github.com/lspaccatrosi16/go-libs/structures/cartesian"
	"github.com/lspaccatrosi16/tmappr/lib/types"
	"github.com/lspaccatrosi16/tmappr/lib/util"
)

func CombineSegments(pathings []*types.PathedLine, maxX, maxY int) (*types.PathedSystem, *cartesian.CoordinateGrid[bool]) {
	logger := logging.GetLogger()
	util.DebugSection("Combine Lines into Greater Map")

	system := types.PathedSystem{}

	for _, path := range pathings {
		logger.Debug(fmt.Sprintf("Handle %s", path.Line.Code))
		for _, section := range path.Segments {
			includeSegment(&system, path.Line, section)
		}
	}

	grid := cartesian.CoordinateGrid[bool]{}

	for x := 0; x <= maxX; x++ {
		for y := 0; y <= maxY; y++ {
			grid.Add(cartesian.Coordinate{x, y}, false)
		}
	}

	for _, seg := range system.Segments {
		for _, point := range seg.Points() {
			grid.Add(point, true)
		}
	}

	return &system, &grid
}

func includeSegment(system *types.PathedSystem, line *types.Line, segment types.LineSegment) {
	existingSegment := system.FindCSegment(segment)
	if existingSegment == nil {
		newSeg := types.CompoundSegment{
			Lines:       []*types.Line{line},
			LineSegment: segment,
		}
		system.AddSegment(&newSeg)
	} else {
		if existingSegment.PointInLine(segment.End) {
			produced := existingSegment.Subsegment(segment.Start, segment.End, line)
			system.AddSegment(produced...)
		} else {
			if segment.Gradient[0] > 0 || segment.Gradient[1] > 0 {
				newSeg := types.LineSegment{
					Start:    existingSegment.End,
					End:      segment.End,
					Gradient: segment.Gradient,
				}
				includeSegment(system, line, newSeg)
			} else {
				newSeg := types.LineSegment{
					Start:    segment.Start,
					End:      existingSegment.Start,
					Gradient: segment.Gradient,
				}
				includeSegment(system, line, newSeg)

			}
		}

	}
}
