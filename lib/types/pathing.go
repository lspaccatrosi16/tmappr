package types

import (
	"fmt"

	"github.com/lspaccatrosi16/go-libs/structures/cartesian"
)

type CompoundSegment struct {
	Lines []*Line
	LineSegment
}

func (c *CompoundSegment) Subsegment(start, end cartesian.Coordinate, newLine *Line) []*CompoundSegment {
	newSegments := []*CompoundSegment{}

	if c.Start != start {
		newEnd := c.Start
		for {
			newEnd = newEnd.Transform(c.Gradient[0], c.Gradient[1])
			if newEnd == start || newEnd == end || newEnd == c.End {
				break
			}
		}
		newSegments = append(newSegments, &CompoundSegment{
			Lines: c.Lines,
			LineSegment: LineSegment{
				Start:    c.Start,
				End:      newEnd,
				Gradient: c.Gradient,
			},
		})
	}

	newSegments = append(newSegments, &CompoundSegment{
		Lines: append(c.Lines, newLine),
		LineSegment: LineSegment{
			Start:    start,
			End:      end,
			Gradient: c.Gradient,
		},
	})

	if c.End != end {
		newStart := c.End
		for {
			newStart = newStart.Transform(-c.Gradient[0], -c.Gradient[1])
			if newStart == start || newStart == end || newStart == c.Start {
				break
			}
		}
		newSegments = append(newSegments, &CompoundSegment{
			Lines: c.Lines,
			LineSegment: LineSegment{
				Start:    newStart,
				End:      c.End,
				Gradient: c.Gradient,
			},
		})
	}

	return newSegments
}

type PathedSystem struct {
	Segments []*CompoundSegment
}

func (p *PathedSystem) FindCSegment(seg LineSegment) *CompoundSegment {
	for _, segment := range p.Segments {
		if segment.PointInLine(seg.Start) && segment.Gradient == seg.Gradient {
			return segment
		}
	}
	return nil

}

func (p *PathedSystem) AddSegment(c ...*CompoundSegment) {
	p.Segments = append(p.Segments, c...)
}

type LineSegment struct {
	Start    cartesian.Coordinate
	End      cartesian.Coordinate
	Gradient [2]int
}

func (l *LineSegment) PointInLine(c cartesian.Coordinate) bool {
	points := l.Points()
	for _, p := range points {
		if p == c {
			return true
		}
	}
	return false
}

func (l *LineSegment) Points() []cartesian.Coordinate {
	points := []cartesian.Coordinate{}
	cp := l.Start
	for {
		points = append(points, cp)
		cp = cp.Transform(l.Gradient[0], l.Gradient[1])
		if cp == l.End {
			break
		}
	}
	points = append(points, l.End)
	return points
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
				Start:    segmentStart,
				End:      prev,
				Gradient: curDirection,
			}
			p.Segments = append(p.Segments, segment)
			curDirection = newDirection
			segmentStart = prev
		}
		prev = path[i]
	}

	finalSegment := LineSegment{
		Start:    segmentStart,
		End:      prev,
		Gradient: curDirection,
	}

	p.Segments = append(p.Segments, finalSegment)

	return nil
}

func cmp_direction(d1, d2 [2]int) bool {
	return d1[0] == d2[0] && d1[1] == d2[1]
}
