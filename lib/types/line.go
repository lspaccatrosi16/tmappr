package types

import (
	"fmt"
	"strings"

	"github.com/lspaccatrosi16/tmappr/lib/util"
)

type Line struct {
	Name   string
	Code   string
	Colour string
	Stops  []*Stop
}

func ParseLine(raw string) (*Line, error) {
	components := strings.Split(util.Trim(raw), ",")
	if len(components) != 3 {
		return nil, fmt.Errorf("expected entry to have 3 components, not %d", len(components))
	}

	for i := 0; i < len(components); i++ {
		components[i] = util.Trim(components[i])
	}

	if !strings.HasPrefix(components[0], "\"") || !strings.HasSuffix(components[0], "\"") {
		return nil, fmt.Errorf("entry name must be double quoted (%s)", components[0])
	}

	if !strings.HasPrefix(components[1], "#") {
		return nil, fmt.Errorf("colour be a hex formatted value starting with #")
	}

	if len(util.Trim(components[1])) != 7 {
		return nil, fmt.Errorf("colour be a hex formatted value starting with #")
	}

	return &Line{
		Name:   components[0][1 : len(components[0])-1],
		Colour: components[1],
		Code:   components[2],
	}, nil
}
