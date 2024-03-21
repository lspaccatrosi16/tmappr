package engine

import (
	"fmt"
	"math"

	"github.com/lspaccatrosi16/go-cli-tools/logging"
	"github.com/lspaccatrosi16/go-libs/structures/cartesian"
	"github.com/lspaccatrosi16/tmappr/lib/types"
	"github.com/lspaccatrosi16/tmappr/lib/util"
)

const Resolution = 2

func RunEngine(config *types.AppConfig, lineMap *map[string]*types.Line, stopMap *map[string]*types.Stop) ([]*types.PathedLine, int, int, error) {
	logger := logging.GetLogger()

	util.DebugSection("Running Pathing Engine")

	grid := cartesian.CoordinateGrid[int]{}
	coordinates := []cartesian.Coordinate{}
	cStopMap := map[cartesian.Coordinate]*types.Stop{}

	for _, s := range *stopMap {
		coord := approxCoordinate(s.Coordinates[0], s.Coordinates[1])
		grid.Add(coord, s.Id)
		coordinates = append(coordinates, coord)
		cStopMap[coord] = s
	}

	maxX, maxY := grid.MaxBounds()

	pathMap := map[string]*types.PathedLine{}
	pathings := []*types.PathedLine{}

	for name, line := range *lineMap {
		logger.Debug(fmt.Sprintf("Pathfind %s", name))

		path, err := GetLinePath(line, &grid, stopMap, &cStopMap)
		if err != nil {
			return nil, 0, 0, err
		}

		logger.Debug(path.String())

		pathMap[name] = path
		pathings = append(pathings, path)

	}
	return pathings, maxX, maxY, nil
}

func approxCoordinate(x, y float64) cartesian.Coordinate {
	xR := math.Floor(x * Resolution)
	yR := math.Floor(y * Resolution)
	return cartesian.Coordinate{int(xR), int(yR)}
}
