package engine

import (
	"fmt"
	"slices"

	"github.com/lspaccatrosi16/go-cli-tools/logging"
	algog "github.com/lspaccatrosi16/go-libs/algorithms/graph"
	structg "github.com/lspaccatrosi16/go-libs/structures/graph"

	"github.com/lspaccatrosi16/go-libs/structures/cartesian"
	"github.com/lspaccatrosi16/tmappr/lib/engine/astar"
	"github.com/lspaccatrosi16/tmappr/lib/types"
	"github.com/lspaccatrosi16/tmappr/lib/util"
)

type pathLeg [2]cartesian.Coordinate

func GetLinePath(config *types.AppConfig, line *types.Line, grid *cartesian.CoordinateGrid[int], cStopMap *map[cartesian.Coordinate]*types.Stop) (*types.PathedLine, error) {
	util.DebugSection(fmt.Sprintf("Pathfind Line %s", line.Code))

	logger := logging.GetLogger()

	legs := []pathLeg{}
	coordinates := []cartesian.Coordinate{}

	approxCoordinate := approxCoordinateMake(float64(config.Simres))

	for _, s := range line.Stops {
		coord := approxCoordinate(s.ExtCoordinates[0], s.ExtCoordinates[1])
		coordinates = append(coordinates, coord)
	}

	for i, c := range coordinates {
		if i+1 < len(coordinates) {
			legs = append(legs, pathLeg{c, coordinates[i+1]})
		}
	}
	pathSections := [][]cartesian.Coordinate{}

	mx, my := grid.MaxBounds()

	switch config.Algorithm {
	case types.Dijkstra:
		genGraph, graphGridMap := grid.CreateGraph(false, []int{}, true)
		for _, leg := range legs {
			result, err := useDijkstra(leg, *graphGridMap, genGraph)
			if err != nil {
				return nil, err
			}
			pathSections = append(pathSections, result)
		}
	case types.AStar:
		currentDirection := cartesian.NoDirection
		forwardsCost := 0.0
		backwardsCost := 0.0

		forwardsSections := [][]cartesian.Coordinate{}
		backwardsSections := [][]cartesian.Coordinate{}
		for _, leg := range legs {
			result, endingDirection, cost := useAstar(leg, currentDirection, mx, my)
			forwardsCost += cost
			currentDirection = endingDirection

			forwardsSections = append(forwardsSections, result)
		}

		slices.Reverse(legs)

		currentDirection = cartesian.NoDirection

		for _, leg := range legs {
			result, endingDirection, cost := useAstar(pathLeg{leg[1], leg[0]}, currentDirection, mx+2, my+2)
			backwardsCost += cost
			currentDirection = endingDirection
			backwardsSections = append(backwardsSections, result)
		}

		if forwardsCost <= backwardsCost {
			pathSections = forwardsSections
		} else {
			pathSections = backwardsSections
		}
	}

	repreGrid := cartesian.CoordinateGrid[string]{}

	for x := 0; x <= mx; x++ {
		for y := 0; y <= my; y++ {
			repreGrid.Add(cartesian.Coordinate{x, y}, " ")
		}
	}

	path := []cartesian.Coordinate{}

	lastEnd := coordinates[0]
	for _, section := range pathSections {
		if len(section) == 0 {
			continue
		}
		if section[0] != lastEnd {
			slices.Reverse(section)
		}

		path = append(path, section...)
		lastEnd = section[len(section)-1]

	}

	parsedPath := []cartesian.Coordinate{}

	for i, p := range path {
		if i > 0 && p == path[i-1] {
			continue
		}

		parsedPath = append(parsedPath, p)
		if _, ok := (*cStopMap)[p]; ok {
			repreGrid.Add(p, "o")
		} else {
			repreGrid.Add(p, "#")
		}
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

func useDijkstra(leg pathLeg, graphGridMap map[cartesian.Coordinate]*cartesian.GraphGridPoint, genGraph *structg.Graph) ([]cartesian.Coordinate, error) {
	path := []cartesian.Coordinate{}
	run, err := algog.RunDijkstra(graphGridMap[leg[0]], graphGridMap[leg[1]], genGraph)
	if err != nil {
		return nil, err
	}

	gpl := new(cartesian.GridPointList).FromGraphNodes(run.DijkstraData.Path)
	for _, gn := range *gpl {
		path = append(path, gn.Point)
	}
	return path, nil
}

func useAstar(leg pathLeg, startingDirection cartesian.Direction, maxX, maxY int) ([]cartesian.Coordinate, cartesian.Direction, float64) {
	path, cost := astar.Pathfind(leg[0], leg[1], startingDirection, maxX, maxY)
	if len(path) < 2 {
		return path, cartesian.NoDirection, cost
	}

	endingDirection := new(cartesian.Direction).FromCoordinates(path[len(path)-2].Subtract(path[len(path)-1]))
	return path, endingDirection, cost
}
