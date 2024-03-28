package engine

import (
	"fmt"

	"github.com/lspaccatrosi16/go-cli-tools/logging"
	"github.com/lspaccatrosi16/go-libs/structures/cartesian"
	"github.com/lspaccatrosi16/tmappr/lib/types"
	"github.com/lspaccatrosi16/tmappr/lib/util"
)

func CombineSegments(config *types.AppConfig, pathings []*types.PathedLine, maxX, maxY int) (*types.PathedSystem, *cartesian.CoordinateGrid[int]) {
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

	refineSystem(config, &system)

	return &system, &grid
}

func findConflicts(system *types.PathedSystem) [][]*types.CompoundSegment {
	used := map[cartesian.Coordinate]map[cartesian.Direction][]*types.CompoundSegment{}

	for _, seg := range system.Segments {
		for _, point := range seg.Points() {
			if used[point] == nil {
				used[point] = map[cartesian.Direction][]*types.CompoundSegment{}
			}

			used[point][seg.Gradient] = append(used[point][seg.Gradient], seg)
		}
	}

	combinedList := [][]*types.CompoundSegment{}

	for _, c := range used {
		for dir, segs := range c {
			if len(segs) > 1 && !compountListContains(combinedList, segs, dir) {
				combinedList = append(combinedList, segs)
			}
		}
	}
	return combinedList
}

var refineCount = 0

func refineSystem(config *types.AppConfig, system *types.PathedSystem) {

	util.DebugSection("Refine System")

	combinedList := findConflicts(system)

	for len(combinedList) > 0 {
		refineCount++

		if refineCount > config.RefineCycles {
			break
		}

		list := combinedList[0]
		if len(list) < 2 {
			continue
		}
		system.RemoveSegment(list[0])
		system.RemoveSegment(list[1])
		produced := joinSegments(list[0], list[1].LineSegment, list[1].Lines)
		system.AddSegment(produced...)
		combinedList = findConflicts(system)
	}

}

func compountListContains(list [][]*types.CompoundSegment, item []*types.CompoundSegment, dir cartesian.Direction) bool {
	seen := map[*types.CompoundSegment]bool{}

	for _, seg := range item {
		seen[seg] = true
	}

	for _, l := range list {
		if len(l) != len(item) {
			continue
		}

		match := true
		for _, seg := range l {
			_, ok := seen[seg]
			if !ok || seg.Gradient != dir {
				match = false
				break
			}

		}

		if match {
			return true
		}
	}
	return false
}

func joinSegments(i1 *types.CompoundSegment, i2 types.LineSegment, lines []*types.Line) []*types.CompoundSegment {
	if i1.PointInLine(i2.Start) && i1.PointInLine(i2.End) {
		// segment is entirely within existing segment
		produced := i1.Subsegment(i2.Start, i2.End, lines...)
		return produced
	} else if i2.PointInLine(i1.Start) && i2.PointInLine(i1.End) {
		// existing segment is entirely within segment
		invSeg := types.NewCompoundSegment(&i2, lines...)
		produced := invSeg.Subsegment(i1.Start, i1.End, i1.Lines...)
		return produced
	} else {
		// segment and existing segment overlap
		if i1.Start == i2.End || i1.End == i2.Start {
			// segment and existing segment are adjacent
			s := types.NewCompoundSegment(&i2, lines...)
			return []*types.CompoundSegment{i1, s}
		} else if i1.PointInLine(i2.Start) {
			// segment starts within existing segment
			// so end is beyond existing segment
			l1 := types.NewLineSegment(i1.Start, i2.Start, i2.Gradient)
			s1 := types.NewCompoundSegment(&l1, i1.Lines...)

			l2 := types.NewLineSegment(i2.Start, s1.End, i1.Gradient)
			s2 := types.NewCompoundSegment(&l2, append(i1.Lines, lines...)...)

			l3 := types.NewLineSegment(s1.End, i2.End, i2.Gradient)
			s3 := types.NewCompoundSegment(&l3, lines...)

			return []*types.CompoundSegment{s1, s2, s3}

		} else if i1.PointInLine(i2.End) {
			// segment ends within existing segment
			// so start is before existing segment

			l1 := types.NewLineSegment(i2.Start, i1.Start, i2.Gradient)
			s1 := types.NewCompoundSegment(&l1, i1.Lines...)

			l2 := types.NewLineSegment(s1.Start, i2.End, i1.Gradient)
			s2 := types.NewCompoundSegment(&l2, append(i1.Lines, lines...)...)

			l3 := types.NewLineSegment(i2.End, s1.End, i2.Gradient)
			s3 := types.NewCompoundSegment(&l3, lines...)

			return []*types.CompoundSegment{s1, s2, s3}
		} else {
			panic("segments overlap but not in a way I can handle")
		}
	}
}

func includeSegment(system *types.PathedSystem, line *types.Line, segment types.LineSegment) {
	var existingSegment *types.CompoundSegment
	forwardsSegment := system.FindCSegment(segment)
	backwardsSegment := system.FindCSegment(segment.Reverse())

	if forwardsSegment != nil {
		existingSegment = forwardsSegment
		system.RemoveSegment(forwardsSegment)
	} else if backwardsSegment != nil {
		existingSegment = backwardsSegment
		system.RemoveSegment(backwardsSegment)
	}

	if existingSegment == nil {
		newSeg := types.NewCompoundSegment(&segment, line)
		system.AddSegment(newSeg)
	} else {
		produced := joinSegments(existingSegment, segment, []*types.Line{line})
		system.AddSegment(produced...)
	}
}
