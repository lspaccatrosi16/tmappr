package types

import (
	"bytes"
	"fmt"
	"math"
	"sort"
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

func (c *CompoundSegment) Subsegment(start, end cartesian.Coordinate, newLines ...*Line) []*CompoundSegment {
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
		l := NewLineSegment(c.Start, newEnd, c.Gradient)
		newSegments = append(newSegments, NewCompoundSegment(&l, c.Lines...))
	}

	l := NewLineSegment(start, end, c.Gradient)

	newSegments = append(newSegments, NewCompoundSegment(&l, append(c.Lines, newLines...)...))

	if c.End != end {
		newStart := c.End
		for {
			newStart = newStart.Add(c.Gradient.Opposite().Coordinates())
			if newStart == start || newStart == end || newStart == c.Start {
				break
			}
		}
		l := NewLineSegment(newStart, c.End, c.Gradient)
		newSegments = append(newSegments, NewCompoundSegment(&l, c.Lines...))
	}

	return newSegments
}

func (c *CompoundSegment) ReorderLines(order []int) {
	newArr := make([]*Line, len(c.Lines))
	newXO := make([]int, len(c.Lines))
	newYO := make([]int, len(c.Lines))

	for i, o := range order {
		newArr[i] = c.Lines[o]
		newXO[i] = c.XOffsets[o]
		newYO[i] = c.YOffsets[o]
	}

	c.Lines = newArr
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
	for _, segment := range p.Segments {
		if segment.PointInLine(seg.Start) && segment.Gradient == seg.Gradient {
			return segment
		}
	}
	return nil
}

func (p *PathedSystem) RemoveSegment(seg *CompoundSegment) {
	for i, segment := range p.Segments {
		if segment == seg {
			p.Segments = append(p.Segments[:i], p.Segments[i+1:]...)
			return
		}
	}
}

type csList []*CompoundSegment

func (c csList) Len() int           { return len(c) }
func (c csList) Swap(i, j int)      { c[i], c[j] = c[j], c[i] }
func (c csList) Less(i, j int) bool { return len(c[i].Lines) < len(c[j].Lines) }

func (p *PathedSystem) FindSegmentsWithPoint(c cartesian.Coordinate) []*CompoundSegment {
	found := csList{}

	for _, seg := range p.Segments {
		if seg.PointInLine(c) {
			found = append(found, seg)
		}
	}

	sort.Sort(found)

	return found
}

func (p *PathedSystem) FindSegmentWithPointLine(c cartesian.Coordinate, l *Line) []*CompoundSegment {
	found := []*CompoundSegment{}
	for _, seg := range p.Segments {
		if seg.PointInLine(c) {
			for _, line := range seg.Lines {
				if line == l {
					found = append(found, seg)
				}
			}
		}
	}
	return found
}

func (p *PathedSystem) AddSegment(c ...*CompoundSegment) {
	p.Segments = append(p.Segments, c...)
}

func NewLineSegment(start, end cartesian.Coordinate, gradient cartesian.Direction) LineSegment {
	if new(cartesian.Direction).FromCoordinates(end.Subtract(start)) != gradient {
		panic("start and end do not match gradient")
	}

	return LineSegment{
		Start:    start,
		End:      end,
		Gradient: gradient,
	}
}

func NewCompoundSegment(ls *LineSegment, lines ...*Line) *CompoundSegment {
	return &CompoundSegment{
		Lines:       lines,
		LineSegment: *ls,
		XOffsets:    make([]int, len(lines)),
		YOffsets:    make([]int, len(lines)),
	}
}

type LineSegment struct {
	Start    cartesian.Coordinate
	End      cartesian.Coordinate
	Gradient cartesian.Direction
}

func (l *LineSegment) Length() float64 {
	dx := l.End[0] - l.Start[0]
	dy := l.End[1] - l.Start[1]
	return math.Sqrt(float64(dx*dx + dy*dy))
}

func (l *LineSegment) MLength() int {
	dx := math.Abs(float64(l.End[0] - l.Start[0]))
	dy := math.Abs(float64(l.End[1] - l.Start[1]))
	return int(dx + dy)
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
