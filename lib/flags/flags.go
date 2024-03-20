package flags

import (
	"flag"
	"fmt"
	"strings"

	"github.com/lspaccatrosi16/tmappr/lib/types"
)

var input = flag.String("i", "", "input location")
var output = flag.String("o", "", "output location")
var format = flag.String("f", "", "output format")
var width = flag.Int("w", 0, "line width")
var lineEnding = flag.String("l", "", "line ending")
var verbose = flag.Bool("v", false, "verbose logging")

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

	return &types.AppConfig{
		Input:     *input,
		Output:    *output,
		Format:    chosenFormat,
		Linewidth: *width,
		Ending:    chosenEnding,
		Verbose:   *verbose,
	}, nil
}
