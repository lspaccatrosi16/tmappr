package types

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/lspaccatrosi16/tmappr/lib/util"
)

type Stop struct {
	Name        string
	Coordinates [2]float64
	Id          int
	Lines       []string
}

func ParseStop(raw string) (*Stop, error) {
	components := strings.Split(util.Trim(raw), ",")

	if len(components) != 2 {
		return nil, fmt.Errorf("expected component to have 2 parts, not %d", len(components))
	}

	for i := 0; i < len(components); i++ {
		components[i] = util.Trim(components[i])
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
	coordinates[0], err = strconv.ParseFloat(util.Trim(coordValues[0]), 64)

	if err != nil {
		return nil, fmt.Errorf("error parsing x coordinate (%s)", components[0])
	}
	coordinates[1], err = strconv.ParseFloat(util.Trim(coordValues[1]), 64)

	if err != nil {
		return nil, fmt.Errorf("error parsing y coordinate (%s)", components[0])
	}

	return &Stop{
		Name:        components[0][1 : len(components[0])-1],
		Coordinates: [2]float64{coordinates[0] + 1, coordinates[1] + 1},
	}, nil
}
