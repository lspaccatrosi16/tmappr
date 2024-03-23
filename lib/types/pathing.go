package types

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/lspaccatrosi16/go-libs/structures/cartesian"
)

type CompoundSegment struct {
	Lines []*Line
	LineSegment
	XOffsets []int
	YOffsets []int
}

func (c *CompoundSegment) String() string {
	lineNames := []string{}
	for _, l := range c.Lines {
		lineNames = append(lineNames, l.Code)
	}
	return fmt.Sprintf("%s: %s", c.LineSegment.String(), strings.Join(lineNames, ", "))
}

func (c *CompoundSegment) Subsegment(start, end cartesian.Coordinate, newLine *Line) []*CompoundSegment {
	newSegments := []*CompoundSegment{}

	if c.Gradient == new(cartesian.Direction).FromCoordinates(end.Subtract(start)).Opposite() {
		start, end = end, start
	}

	if c.Start != start {
		newEnd := c.Start
		for {
			newEnd = newEnd.Add(c.Gradient.Coordinates())
			if newEnd == start || newEnd == end || newEnd == c.End {
				break
			}
		}
		newSegments = append(newSegments, &CompoundSegment{
			Lines:       c.Lines,
			LineSegment: NewLineSegment(c.Start, newEnd, c.Gradient),
		})
	}

	newSegments = append(newSegments, &CompoundSegment{
		Lines:       append(c.Lines, newLine),
		LineSegment: NewLineSegment(start, end, c.Gradient),
	})

	if c.End != end {
		newStart := c.End
		for {
			newStart = newStart.Add(c.Gradient.Opposite().Coordinates())
			if newStart == start || newStart == end || newStart == c.Start {
				break
			}
		}
		newSegments = append(newSegments, &CompoundSegment{
			Lines:       c.Lines,
			LineSegment: NewLineSegment(newStart, c.End, c.Gradient),
		})
	}

	return newSegments
}

type PathedSystem struct {
	Segments []*CompoundSegment
}

func (p *PathedSystem) String() string {
	buf := bytes.NewBuffer(nil)
	for _, s := range p.Segments {
		fmt.Fprintln(buf, s.String())
	}
	return buf.String()
}

func (p *PathedSystem) FindCSegment(seg LineSegment) *CompoundSegment {
	for i, segment := range p.Segments {
		if segment.PointInLine(seg.Start) && segment.Gradient == seg.Gradient {
			p.Segments = append(p.Segments[:i], p.Segments[i+1:]...)
			return segment
		}
	}
	return nil
}

func (p *PathedSystem) FindPrimarySegmentWithPoint(c cartesian.Coordinate) *CompoundSegment {
	var chosen *CompoundSegment
	maxLines := 0
	for _, seg := range p.Segments {
		if seg.PointInLine(c) {
			if len(seg.Lines) > maxLines {
				chosen = seg
				maxLines = len(seg.Lines)
			}
		}
	}

	return chosen
}

func (p *PathedSystem) AddSegment(c ...*CompoundSegment) {
	p.Segments = append(p.Segments, c...)
}

func NewLineSegment(start, end cartesian.Coordinate, gradient cartesian.Direction) LineSegment {
	return LineSegment{
		Start:    start,
		End:      end,
		Gradient: gradient,
	}
}

type LineSegment struct {
	Start    cartesian.Coordinate
	End      cartesian.Coordinate
	Gradient cartesian.Direction
}

func (l *LineSegment) Reverse() LineSegment {
	return NewLineSegment(l.End, l.Start, l.Gradient.Opposite())
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

	if l.Start == l.End {
		return []cartesian.Coordinate{l.Start}
	}

	for {
		points = append(points, cp)
		cp = cp.TransformInDirection(l.Gradient)
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

	curDirection := new(cartesian.Direction).FromCoordinates(path[1].Subtract(prev))

	for i := 1; i < len(path); i++ {
		newDirection := new(cartesian.Direction).FromCoordinates(path[i].Subtract(prev))
		if curDirection != newDirection {
			segment := NewLineSegment(segmentStart, prev, curDirection)

			p.Segments = append(p.Segments, segment)
			curDirection = newDirection
			segmentStart = prev
		}
		prev = path[i]
	}

	finalSegment := NewLineSegment(segmentStart, prev, curDirection)

	p.Segments = append(p.Segments, finalSegment)

	return nil
}
