package types

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/lspaccatrosi16/go-libs/structures/cartesian"
	"github.com/lspaccatrosi16/go-libs/structures/enum"
	"github.com/lspaccatrosi16/tmappr/lib/util"
)

type Stop struct {
	Name        string
	Coordinates [2]float64
	Id          int
	Lines       []string
	Type        StopType
	Position    cartesian.Direction
}

func parseCoordinates(i, name string) ([2]int, error) {
	if !strings.HasPrefix(i, "[") || !strings.HasSuffix(i, "]") {
		return [2]int{}, fmt.Errorf("coordinates must be surrounded by brackets (%s)", name)
	}

	coordValues := strings.Split(i[1:len(i)-1], "/")
	if len(coordValues) != 2 {
		return [2]int{}, fmt.Errorf("coordinates must have 2 components (%s)", name)
	}

	var err error

	coordinates := [2]int{}
	coordinates[0], err = strconv.Atoi(util.Trim(coordValues[0]))

	if err != nil {
		return [2]int{}, fmt.Errorf("error parsing x coordinate (%s)", name)
	}
	coordinates[1], err = strconv.Atoi(util.Trim(coordValues[1]))

	if err != nil {
		return [2]int{}, fmt.Errorf("error parsing y coordinate (%s)", name)
	}

	return coordinates, nil

}

func ParseStop(raw string) (*Stop, error) {
	components := strings.Split(util.Trim(raw), ",")

	if len(components) < 2 {
		return nil, fmt.Errorf("expected component to have at least 2 parts, not %d", len(components))
	}

	for i := 0; i < len(components); i++ {
		components[i] = util.Trim(components[i])
	}

	if !strings.HasPrefix(components[0], "\"") || !strings.HasSuffix(components[0], "\"") {
		return nil, fmt.Errorf("entry name must be double quoted (%s)", components[0])
	}

	coordinates, err := parseCoordinates(components[1], components[0])
	if err != nil {
		return nil, err
	}

	var stopType = AutoStopType

	if len(components) >= 3 && components[2] != "" {
		stopType = GetStopType(components[2])
		if !stopType.IsValid() {
			return nil, fmt.Errorf("invalid stop type, %s\nAvailable formats: %s", components[2], strings.Join(enum.All[StopType](), ", "))
		}
	}

	var stopPosition = cartesian.NoDirection

	if len(components) >= 4 && components[3] != "" {
		stopPosition = GetStopPosition(components[3])
		if !stopPosition.IsValid() {
			return nil, fmt.Errorf("invalid stop position, %s\nAvailable formats: %s", components[3], strings.Join(enum.All[cartesian.Direction](), ", "))
		}
	}

	return &Stop{
		Name:        components[0][1 : len(components[0])-1],
		Coordinates: [2]float64{float64(coordinates[0]), float64(coordinates[1])},
		Type:        stopType,
		Position:    stopPosition,
	}, nil
}
