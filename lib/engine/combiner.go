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
	var existingSegment *types.CompoundSegment
	forwardsSegment := system.FindCSegment(segment)
	backwardsSegment := system.FindCSegment(segment.Reverse())

	if forwardsSegment != nil {
		existingSegment = system.FindCSegmentRemove(segment)
	} else if backwardsSegment != nil {
		existingSegment = system.FindCSegmentRemove(segment.Reverse())
	}

	if existingSegment == nil {
		newSeg := types.CompoundSegment{
			Lines:       []*types.Line{line},
			LineSegment: segment,
		}
		system.AddSegment(&newSeg)
	} else {
		if existingSegment.PointInLine(segment.Start) && existingSegment.PointInLine(segment.End) {
			// segment is entirely within existing segment
			produced := existingSegment.Subsegment(segment.Start, segment.End, line)
			system.AddSegment(produced...)
		} else if segment.PointInLine(existingSegment.Start) && segment.PointInLine(existingSegment.End) {
			// existing segment is entirely within segment
			invSeg := types.CompoundSegment{
				Lines:       []*types.Line{line},
				LineSegment: segment,
			}
			produced := invSeg.Subsegment(existingSegment.Start, existingSegment.End, existingSegment.Lines...)
			system.AddSegment(produced...)
		} else {
			// segment and existing segment overlap
			if existingSegment.Start == segment.End || existingSegment.End == segment.Start {
				// segment and existing segment are adjacent
				s := types.CompoundSegment{
					LineSegment: segment,
					Lines:       []*types.Line{line},
				}
				system.AddSegment(existingSegment, &s)
			} else if existingSegment.PointInLine(segment.Start) {
				// segment starts within existing segment
				// so end is beyond existing segment
				s1 := types.CompoundSegment{
					LineSegment: types.LineSegment{
						Start: existingSegment.Start,
						End:   segment.Start,
					},
					Lines: existingSegment.Lines,
				}

				s2 := types.CompoundSegment{
					LineSegment: types.LineSegment{
						Start: segment.Start,
						End:   existingSegment.End,
					},
					Lines: append(existingSegment.Lines, line),
				}
				s3 := types.CompoundSegment{
					LineSegment: types.LineSegment{
						Start: existingSegment.End,
						End:   segment.End,
					},
					Lines: []*types.Line{line},
				}

				system.AddSegment(&s1, &s2, &s3)

			} else if existingSegment.PointInLine(segment.End) {
				// segment ends within existing segment
				// so start is before existing segment

				s1 := types.CompoundSegment{
					LineSegment: types.LineSegment{
						Start: segment.Start,
						End:   existingSegment.Start,
					},
					Lines: []*types.Line{line},
				}

				s2 := types.CompoundSegment{
					LineSegment: types.LineSegment{
						Start: existingSegment.Start,
						End:   segment.End,
					},
					Lines: append(existingSegment.Lines, line),
				}
				s3 := types.CompoundSegment{
					LineSegment: types.LineSegment{
						Start: segment.End,
						End:   existingSegment.End,
					},
					Lines: existingSegment.Lines,
				}

				system.AddSegment(&s1, &s2, &s3)
			} else {
				panic("segments overlap but not in a way I can handle")
			}
		}
	}
}
