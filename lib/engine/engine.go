package engine

import (
	"math"

	"github.com/lspaccatrosi16/go-libs/structures/cartesian"
	"github.com/lspaccatrosi16/tmappr/lib/types"
)

const Resolution = 2

func RunEngine(config *types.AppConfig, lineMap *map[string]*types.Line, stopMap *map[string]*types.Stop, ratios [2]float64) error {
	grid := cartesian.CoordinateGrid[int]{}
	coordinates := []cartesian.Coordinate{}
	cStopMap := map[cartesian.Coordinate]*types.Stop{}

	for _, s := range *stopMap {
		coord := approxCoordinate(s.Coordinates[0], s.Coordinates[1])
		grid.Add(coord, s.Id)
		coordinates = append(coordinates, coord)
		cStopMap[coord] = s
	}

	for _, line := range *lineMap {
		err := GetLinePath(line, &grid, stopMap, &cStopMap)
		if err != nil {
			return err
		}
	}
	return nil
}

func approxCoordinate(x, y float64) cartesian.Coordinate {
	xR := math.Floor(x * Resolution)
	yR := math.Floor(y * Resolution)
	return cartesian.Coordinate{int(xR), int(yR)}
}
