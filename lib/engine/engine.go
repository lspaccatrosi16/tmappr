package engine

import (
	"fmt"
	"math"

	"github.com/lspaccatrosi16/go-cli-tools/logging"
	"github.com/lspaccatrosi16/go-libs/structures/cartesian"
	"github.com/lspaccatrosi16/tmappr/lib/types"
	"github.com/lspaccatrosi16/tmappr/lib/util"
)

const BorderCoordOffet = 1

func RunEngine(config *types.AppConfig, data *types.AppData) error {
	logger := logging.GetLogger()

	util.DebugSection("Running Pathing Engine")

	approxCoordinate := approxCoordinateMake(float64(config.Simres))

	grid := cartesian.CoordinateGrid[int]{}
	cStopMap := map[cartesian.Coordinate]*types.Stop{}

	for _, s := range data.Stops {
		coord := approxCoordinate(s.ExtCoordinates[0], s.ExtCoordinates[1])
		grid.Add(coord, s.Id)
		cStopMap[coord] = s
		s.IntCoordinates = coord
	}

	maxX, maxY := grid.MaxBounds()

	pathings := []*types.PathedLine{}

	for _, line := range data.LinesNames {
		logger.Debug(fmt.Sprintf("Pathfind %s", line.Code))

		path, err := GetLinePath(config, line, &grid, &cStopMap)
		if err != nil {
			return err
		}

		logger.Debug(path.String())

		pathings = append(pathings, path)

	}

	combined, combinedGrid := CombineSegments(config, pathings, maxX+1, maxY+1)

	data.Pathings = combined
	data.UsedGrid = combinedGrid
	data.CStopMap = cStopMap
	data.MaxX = maxX + BorderCoordOffet
	data.MaxY = maxY + BorderCoordOffet

	return nil
}

func approxCoordinateMake(res float64) func(x, y float64) cartesian.Coordinate {
	return func(x, y float64) cartesian.Coordinate {
		xR := math.Floor(x*res) + BorderCoordOffet
		yR := math.Floor(y*res) + BorderCoordOffet
		return cartesian.Coordinate{int(xR), int(yR)}

	}
}
