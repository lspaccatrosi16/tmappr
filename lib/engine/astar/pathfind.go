package astar

import (
	"math"

	"github.com/lspaccatrosi16/go-libs/structures/cartesian"
	"github.com/lspaccatrosi16/go-libs/structures/mpq"
)

const directionCost = 1

type Node struct {
	cartesian.Coordinate
	Heuristic float64
	Cost      float64
	Direction cartesian.Direction
}

func (n *Node) IWeight() int {
	return int(math.Ceil(n.Heuristic + n.Cost))
}

func (n *Node) Weight() float64 {
	return n.Heuristic + n.Cost
}

func neighbors(node *Node, currentDirection cartesian.Direction) []*Node {
	neighbors := []*Node{}
	for _, direction := range cartesian.CardinalPositions() {
		if direction == currentDirection.Opposite() {
			continue
		}
		neighbor := node.Add(direction.Coordinates())

		if neighbor[0] < 0 || neighbor[1] < 0 {
			continue
		}

		neighbors = append(neighbors, &Node{Coordinate: neighbor, Direction: direction})
	}

	return neighbors
}

func heuristic(a, b cartesian.Coordinate) float64 {
	dx := math.Abs(float64(b[0] - a[0]))
	dy := math.Abs(float64(b[1] - a[1]))

	if dx == dy {
		return dx
	} else if dx > dy {
		return dx
	} else {
		return dy
	}
}

func Pathfind(start, end cartesian.Coordinate, startingDirection cartesian.Direction) ([]cartesian.Coordinate, float64) {
	startNode := &Node{Coordinate: start, Direction: startingDirection}
	endNode := &Node{Coordinate: end}

	bestScore := map[cartesian.Coordinate]float64{}
	prev := map[*Node]*Node{}

	openSet := mpq.Queue[*Node]{}
	openSet.Add(startNode, 0)

	for openSet.Len() > 0 {
		current := openSet.Pop()
		if current.Coordinate == endNode.Coordinate {
			path := []cartesian.Coordinate{}
			for current != nil {
				path = append(path, current.Coordinate)
				current = prev[current]
			}
			return path, bestScore[endNode.Coordinate]
		}

		for _, neighbor := range neighbors(current, current.Direction) {
			neighbor.Heuristic = heuristic(neighbor.Coordinate, endNode.Coordinate)
			neighbor.Cost = current.Cost + 1

			if current.Direction != cartesian.NoDirection {
				if current.Direction != neighbor.Direction {
					cwDiff := current.Direction.NumberCw(neighbor.Direction)
					acwDiff := current.Direction.NumberAcw(neighbor.Direction)
					if cwDiff < 2 || acwDiff < 2 {
						neighbor.Cost += directionCost
					} else if cwDiff == 2 || acwDiff == 2 {
						neighbor.Cost += 1.5 * directionCost
					} else {
						neighbor.Cost += 2 * directionCost
					}
				}
			}

			pathWeight := neighbor.Weight()

			best, ok := bestScore[neighbor.Coordinate]
			if !ok {
				openSet.Add(neighbor, neighbor.IWeight())
			} else {
				if best <= pathWeight {
					continue
				}
			}
			bestScore[neighbor.Coordinate] = pathWeight
			prev[neighbor] = current
		}
	}
	return nil, 0
}
