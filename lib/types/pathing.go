package types

import (
	"fmt"

	"github.com/lspaccatrosi16/go-libs/structures/cartesian"
)

type CompoundSegment struct {
	Start cartesian.Coordinate
	End   cartesian.Coordinate
	Lines []*Line
}

type PathedSystem struct {
	Segments []CompoundSegment
}

type LineSegment struct {
	Start cartesian.Coordinate
	End   cartesian.Coordinate
}

func (l *LineSegment) String() string {
	return fmt.Sprintf("%s => %s", l.Start.String(), l.End.String())
}

type PathedLine struct {
	Segments []LineSegment
	Line     *Line
}

func (p *PathedLine) String() string {
	s := ""
	for _, segment := range p.Segments {
		s += segment.String() + "\n"
	}
	return s
}

func (p *PathedLine) CreateSegments(path []cartesian.Coordinate) error {
	if len(path) < 2 {
		return fmt.Errorf("path is too short")
	}

	prev := path[0]
	segmentStart := path[0]
	curDirection := [2]int{path[1][0] - prev[0], path[1][1] - prev[1]}

	for i := 1; i < len(path); i++ {
		newDirection := [2]int{path[i][0] - prev[0], path[i][1] - prev[1]}
		if !cmp_direction(curDirection, newDirection) {
			segment := LineSegment{
				Start: segmentStart,
				End:   prev,
			}
			p.Segments = append(p.Segments, segment)
			curDirection = newDirection
			segmentStart = prev
		}
		prev = path[i]
	}

	finalSegment := LineSegment{
		Start: segmentStart,
		End:   prev,
	}

	p.Segments = append(p.Segments, finalSegment)

	return nil
}

func cmp_direction(d1, d2 [2]int) bool {
	return d1[0] == d2[0] && d1[1] == d2[1]
}
