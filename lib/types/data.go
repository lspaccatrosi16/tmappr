package types

import "github.com/lspaccatrosi16/go-libs/structures/cartesian"

type AppData struct {
	Stops       []*Stop
	StopNames   map[string]*Stop
	Lines       []*Line
	LinesNames  map[string]*Line
	StopLineMap map[*Stop][]*Line
	CStopMap    map[cartesian.Coordinate]*Stop
	MaxX        int
	MaxY        int
	Pathings    *PathedSystem
	UsedGrid    *cartesian.CoordinateGrid[int]
}
