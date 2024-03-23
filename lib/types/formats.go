package types

import (
	"github.com/lspaccatrosi16/go-libs/structures/cartesian"
	"github.com/lspaccatrosi16/go-libs/structures/enum"
)

type Format int

const (
	SVG Format = iota
	PNG
)

func (f Format) String() string {
	switch f {
	case SVG:
		return "svg"
	case PNG:
		return "png"
	default:
		return enum.InvalidString
	}
}

func (f Format) FromI(i int) enum.Enum {
	return Format(i)
}

func (f Format) IsValid() bool {
	return f.String() != enum.InvalidString
}

type LineEnding int

const (
	LF LineEnding = iota
	CRLF
)

func (l LineEnding) String() string {
	switch l {
	case LF:
		return "LF"
	case CRLF:
		return "CRLF"
	default:
		return enum.InvalidString
	}
}

func (l LineEnding) FromI(i int) enum.Enum {
	return LineEnding(i)
}

func (l LineEnding) IsValid() bool {
	return l.String() != enum.InvalidString
}

type StopType int

const (
	AutoStopType StopType = iota
	Normal
	Interchange
)

func (s StopType) FromI(i int) enum.Enum {
	return StopType(i)
}

func (s StopType) IsValid() bool {
	return s.String() != enum.InvalidString
}

func (s StopType) String() string {
	switch s {
	case Normal:
		return "Normal"
	case Interchange:
		return "Interchange"
	case AutoStopType:
		return "Auto"
	default:
		return enum.InvalidString
	}
}

type Algorithm int

const (
	AStar Algorithm = iota
	Dijkstra
)

func (a Algorithm) String() string {
	switch a {
	case AStar:
		return "A*"
	case Dijkstra:
		return "Dijkstra"
	default:
		return enum.InvalidString
	}
}

func (a Algorithm) FromI(i int) enum.Enum {
	return Algorithm(i)
}

func (a Algorithm) IsValid() bool {
	return a.String() != enum.InvalidString
}

func GetFormat(i string) Format {
	switch i {
	case "svg", "SVG":
		return SVG
	case "png", "PNG":
		return PNG
	default:
		return Format(-1)
	}
}

func GetEnding(i string) LineEnding {
	switch i {
	case "crlf", "CRLF":
		return CRLF
	case "lf", "LF":
		return LF
	default:
		return LineEnding(-1)
	}
}

func GetStopType(i string) StopType {
	switch i {
	case "auto", "Auto", "AUTO":
		return AutoStopType
	case "normal", "Normal", "NORMAL":
		return Normal
	case "interchange", "Interchange", "INTERCHANGE":
		return Interchange
	default:
		return StopType(-1)
	}
}

func GetAlgorithm(i string) Algorithm {
	switch i {
	case "astar", "AStar", "ASTAR":
		return AStar
	case "dijkstra", "Dijkstra", "DIJKSTRA":
		return Dijkstra
	default:
		return Algorithm(-1)
	}
}

func GetStopPosition(i string) cartesian.Direction {
	switch i {
	case "north", "North", "NORTH":
		return cartesian.North
	case "northeast", "NorthEast", "NORTHEAST":
		return cartesian.NorthEast
	case "east", "East", "EAST":
		return cartesian.East
	case "southeast", "SouthEast", "SOUTHEAST":
		return cartesian.SouthEast
	case "south", "South", "SOUTH":
		return cartesian.South
	case "southwest", "SouthWest", "SOUTHWEST":
		return cartesian.SouthWest
	case "west", "West", "WEST":
		return cartesian.West
	case "northwest", "NorthWest", "NORTHWEST":
		return cartesian.NorthWest
	case "auto", "Auto", "AUTO":
		return cartesian.NoDirection
	default:
		return cartesian.Direction(-1)
	}
}
