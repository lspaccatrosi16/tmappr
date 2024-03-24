package engine

import (
	"fmt"

	"github.com/lspaccatrosi16/go-cli-tools/logging"
	"github.com/lspaccatrosi16/go-libs/structures/cartesian"
	"github.com/lspaccatrosi16/tmappr/lib/types"
	"github.com/lspaccatrosi16/tmappr/lib/util"
)

func CombineSegments(pathings []*types.PathedLine, maxX, maxY int) (*types.PathedSystem, *cartesian.CoordinateGrid[int]) {
	logger := logging.GetLogger()
	util.DebugSection("Combine Lines into Greater Map")

	system := types.PathedSystem{}

	for _, path := range pathings {
		logger.Debug(fmt.Sprintf("Handle %s", path.Line.Code))
		for _, section := range path.Segments {
			includeSegment(&system, path.Line, section)
		}
	}

	grid := cartesian.CoordinateGrid[int]{}

	for x := 0; x <= maxX; x++ {
		for y := 0; y <= maxY; y++ {
			grid.Add(cartesian.Coordinate{x, y}, 0)
		}
	}

	for _, seg := range system.Segments {
		for _, point := range seg.Points() {
			grid.Add(point, 1)
		}
	}

	logger.Debug("Used Grid")
	logger.Debug(grid.String())

	return &system, &grid
}

func includeSegment(system *types.PathedSystem, line *types.Line, segment types.LineSegment) {
	var existingSegment *types.Corridor
	forwardsSegment := system.FindCSegment(segment)
	backwardsSegment := system.FindCSegment(segment.Reverse())

	if forwardsSegment != nil {
		existingSegment = forwardsSegment
	} else if backwardsSegment != nil {
		existingSegment = backwardsSegment
	}

	if existingSegment == nil {
		newSeg := types.Corridor{
			Lines:       []*types.Line{line},
			LineSegment: segment,
		}
		system.AddSegment(&newSeg)
	} else {
		if existingSegment.PointInLine(segment.End) {
			produced := existingSegment.Subsegment(segment.Start, segment.End, line)
			system.AddSegment(produced...)
		} else {
			if segment.Gradient.Coordinates()[0] > 0 || segment.Gradient.Coordinates()[1] > 0 {
				newSeg := types.NewLineSegment(existingSegment.End, segment.End, segment.Gradient)
				includeSegment(system, line, newSeg)
			} else {
				newSeg := types.NewLineSegment(segment.Start, existingSegment.Start, segment.Gradient)
				includeSegment(system, line, newSeg)

			}
		}

	}
}
