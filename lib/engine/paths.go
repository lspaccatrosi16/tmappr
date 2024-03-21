package engine

import (
	"fmt"

	"github.com/lspaccatrosi16/go-cli-tools/logging"
	"github.com/lspaccatrosi16/go-libs/algorithms/graph"
	"github.com/lspaccatrosi16/go-libs/structures/cartesian"
	"github.com/lspaccatrosi16/tmappr/lib/types"
	"github.com/lspaccatrosi16/tmappr/lib/util"
)

type pathLeg [2]cartesian.Coordinate

func GetLinePath(line *types.Line, grid *cartesian.CoordinateGrid[int], stopMap *map[string]*types.Stop, cStopMap *map[cartesian.Coordinate]*types.Stop) (*types.PathedLine, error) {
	util.DebugSection(fmt.Sprintf("Pathfind Line %s", line.Code))

	logger := logging.GetLogger()

	legs := []pathLeg{}
	coordinates := []cartesian.Coordinate{}

	for _, s := range line.Stops {
		coord := approxCoordinate((*stopMap)[s].Coordinates[0], (*stopMap)[s].Coordinates[1])
		coordinates = append(coordinates, coord)
	}

	for i, c := range coordinates {
		if i+1 < len(coordinates) {
			legs = append(legs, pathLeg{c, coordinates[i+1]})
		}
	}

	genGraph, graphGridMap := grid.CreateGraph(false, []int{}, true)

	path := []cartesian.Coordinate{}

	for _, leg := range legs {
		run, err := graph.RunDijkstra((*graphGridMap)[leg[0]], (*graphGridMap)[leg[1]], genGraph)
		if err != nil {
			return nil, err
		}

		gpl := new(cartesian.GridPointList).FromGraphNodes(run.DijkstraData.Path)
		for _, gn := range *gpl {
			path = append(path, gn.Point)
		}
	}

	repreGrid := cartesian.CoordinateGrid[string]{}

	mx, my := grid.MaxBounds()

	for x := 0; x <= mx; x++ {
		for y := 0; y <= my; y++ {
			repreGrid.Add(cartesian.Coordinate{x, y}, " ")
		}
	}

	parsedPath := []cartesian.Coordinate{}

	for i, p := range path {
		if i+1 < len(path) && p == path[i+1] {
			continue
		}

		if _, ok := (*cStopMap)[p]; ok {
			repreGrid.Add(p, "o")
		} else {
			repreGrid.Add(p, "#")
		}

		parsedPath = append(parsedPath, p)
	}

	logger.Debug(repreGrid.String())

	pathing := types.PathedLine{
		Line: line,
	}

	err := pathing.CreateSegments(parsedPath)

	if err != nil {
		return nil, err
	}

	return &pathing, nil
}
