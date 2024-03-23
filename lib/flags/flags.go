package flags

import (
	"flag"
	"fmt"
	"strconv"
	"strings"

	"github.com/lspaccatrosi16/go-libs/structures/enum"
	"github.com/lspaccatrosi16/tmappr/lib/types"
)

var input = flag.String("i", "", "input location")
var output = flag.String("o", "", "output location")
var format = flag.String("f", "svg", "output format")
var width = flag.Int("w", 6, "line width")
var lineEnding = flag.String("l", "LF", "line ending")
var verbose = flag.Bool("v", false, "verbose logging")
var res = flag.String("r", "", "resolution (x:y)")
var simRes = flag.Int("sr", 4, "simulation resolution")
var pathfinder = flag.String("p", "astar", "pathfinding algorithm")

func GetFlagData() (*types.AppConfig, error) {
	flag.Parse()

	if *input == "" {
		return nil, fmt.Errorf("missing or empty input location (-i)")
	}

	if *output == "" {
		return nil, fmt.Errorf("missing or empty output location (-o)")
	}

	if *format == "" {
		return nil, fmt.Errorf("missing or empty output format (-f)")
	}

	chosenFormat := types.GetFormat(*format)

	if !chosenFormat.IsValid() {
		return nil, fmt.Errorf("invalid format, %s\nAvailable formats: %s", *format, strings.Join(enum.All[types.Format](), ", "))
	}

	if *width == 0 {
		return nil, fmt.Errorf("missing or empty line width (-w)")
	}

	if *lineEnding == "" {
		return nil, fmt.Errorf("missing or empty line ending (-l)")
	}

	chosenEnding := types.GetEnding(*lineEnding)

	if !chosenEnding.IsValid() {
		return nil, fmt.Errorf("invalid format, %s\nAvailable formats: %s", *format, strings.Join(enum.All[types.LineEnding](), ", "))
	}

	if *res == "" {
		return nil, fmt.Errorf("missing or empty resolution (-r)")
	}

	resComps := strings.Split(*res, ":")
	if len(resComps) != 2 {
		return nil, fmt.Errorf("invalid resolution format (x:y)")
	}

	xRes, err := strconv.Atoi(resComps[0])
	if err != nil {
		return nil, fmt.Errorf("resolution x has invalid format")
	}
	yRes, err := strconv.Atoi(resComps[1])
	if err != nil {
		return nil, fmt.Errorf("resolution y has invalid format")
	}

	if *pathfinder == "" {
		return nil, fmt.Errorf("missing or empty pathfinder (-p)")
	}

	chosenAlgorithm := types.GetAlgorithm(*pathfinder)

	if !chosenAlgorithm.IsValid() {
		return nil, fmt.Errorf("invalid pathfinder, %s\nAvailable pathfinders: %s", *pathfinder, strings.Join(enum.All[types.Algorithm](), ", "))
	}

	return &types.AppConfig{
		Input:     *input,
		Output:    *output,
		Format:    chosenFormat,
		Linewidth: *width,
		Ending:    chosenEnding,
		Verbose:   *verbose,
		XRes:      xRes,
		YRes:      yRes,
		Simres:    *simRes,
		Algorithm: chosenAlgorithm,
	}, nil
}
