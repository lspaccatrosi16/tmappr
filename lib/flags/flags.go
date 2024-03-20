package flags

import (
	"flag"
	"fmt"
	"strconv"
	"strings"

	"github.com/lspaccatrosi16/tmappr/lib/types"
)

var input = flag.String("i", "", "input location")
var output = flag.String("o", "", "output location")
var format = flag.String("f", "", "output format")
var width = flag.Int("w", 0, "line width")
var lineEnding = flag.String("l", "", "line ending")
var verbose = flag.Bool("v", false, "verbose logging")
var res = flag.String("r", "", "resolution (x:y)")

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
		return nil, fmt.Errorf("invalid format, %s\nAvailable formats: %s", *format, strings.Join(types.All[types.Format](), ", "))
	}

	if *width == 0 {
		return nil, fmt.Errorf("missing or empty line width (-w)")
	}

	if *lineEnding == "" {
		return nil, fmt.Errorf("missing or empty line ending (-l)")
	}

	chosenEnding := types.GetEnding(*lineEnding)

	if !chosenEnding.IsValid() {
		return nil, fmt.Errorf("invalid format, %s\nAvailable formats: %s", *format, strings.Join(types.All[types.LineEnding](), ", "))
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

	return &types.AppConfig{
		Input:     *input,
		Output:    *output,
		Format:    chosenFormat,
		Linewidth: *width,
		Ending:    chosenEnding,
		Verbose:   *verbose,
		XRes:      xRes,
		YRes:      yRes,
	}, nil
}
