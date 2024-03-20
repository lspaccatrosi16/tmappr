package engine

import (
	"fmt"
	"math"

	"github.com/lspaccatrosi16/go-libs/algorithms/graph"
	"github.com/lspaccatrosi16/go-libs/structures/cartesian"
	"github.com/lspaccatrosi16/tmappr/lib/types"
)

func ArrangePaths(config *types.AppConfig, lineMap map[string]*types.Line, stops []*types.Stop, ratios [2]float64) error {
	grid := cartesian.CoordinateGrid[int]{}
	coordinates := []cartesian.Coordinate{}

	for i, stop := range stops {
		coord := approxCoordinate(stop.Coordinates[0]*ratios[0], stop.Coordinates[1]*ratios[1])
		grid.Add(coord, i+1)
		coordinates = append(coordinates, coord)
	}

	genGraph, graphGridMap := grid.CreateGraph(false, []int{}, true)

	run, err := graph.RunDijkstra((*graphGridMap)[coordinates[0]], (*graphGridMap)[coordinates[1]], genGraph)
	if err != nil {
		return err
	}

	repr := grid.GraphSearchRepresentation(run)
	fmt.Println(repr)
	return nil
}

func approxCoordinate(x, y float64) cartesian.Coordinate {
	xR := math.Floor(x)
	yR := math.Floor(y)
	return cartesian.Coordinate{int(xR), int(yR)}
}
