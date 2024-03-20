package types

import (
	"fmt"
	"strconv"
	"strings"
)

type Stop struct {
	Name        string
	Coordinates [2]float64
	Lines       []string
}

func ParseStop(raw string) (*Stop, error) {
	components := strings.Split(trim(raw), ",")

	if len(components) < 3 {
		return nil, fmt.Errorf("expected component to have at least 3 parts, not %d", len(components))
	}

	for i := 0; i < len(components); i++ {
		components[i] = trim(components[i])
	}

	if !strings.HasPrefix(components[0], "\"") || !strings.HasSuffix(components[0], "\"") {
		return nil, fmt.Errorf("entry name must be double quoted (%s)", components[0])
	}

	if !strings.HasPrefix(components[1], "[") || !strings.HasSuffix(components[1], "]") {
		return nil, fmt.Errorf("station coordinates must be surrounded by brackets (%s)", components[0])
	}

	coordValues := strings.Split(components[1][1:len(components[1])-1], "/")
	if len(coordValues) != 2 {
		return nil, fmt.Errorf("station coordinates must have 2 components (%s)", components[0])
	}

	var err error

	coordinates := [2]float64{}
	coordinates[0], err = strconv.ParseFloat(trim(coordValues[0]), 64)

	if err != nil {
		return nil, fmt.Errorf("error parsing x coordinate (%s)", components[0])
	}
	coordinates[1], err = strconv.ParseFloat(trim(coordValues[1]), 64)

	if err != nil {
		return nil, fmt.Errorf("error parsing y coordinate (%s)", components[0])
	}

	return &Stop{
		Name:        components[0][1 : len(components[0])-1],
		Coordinates: coordinates,
		Lines:       components[2:],
	}, nil
}
