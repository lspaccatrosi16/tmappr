package types

import (
	"fmt"
	"strings"
)

type Line struct {
	Name   string
	Code   string
	Colour string
}

func trim(s string) string {
	return strings.Trim(s, " \t\r\n")
}

func ParseLine(raw string) (*Line, error) {
	components := strings.Split(trim(raw), ",")
	if len(components) != 3 {
		return nil, fmt.Errorf("expected entry to have 3 components, not %d", len(components))
	}

	for i := 0; i < len(components); i++ {
		components[i] = trim(components[i])
	}

	if !strings.HasPrefix(components[0], "\"") || !strings.HasSuffix(components[0], "\"") {
		return nil, fmt.Errorf("entry name must be double quoted (%s)", components[0])
	}

	if !strings.HasPrefix(components[1], "#") {
		return nil, fmt.Errorf("colour be a hex formatted value starting with #")
	}

	if len(trim(components[1])) != 7 {
		return nil, fmt.Errorf("colour be a hex formatted value starting with #")
	}

	return &Line{
		Name:   components[0][1 : len(components)-1],
		Colour: components[1],
		Code:   components[2],
	}, nil
}
