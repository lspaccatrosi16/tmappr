package io

import (
	"fmt"
	"strings"

	"github.com/lspaccatrosi16/go-cli-tools/logging"
	"github.com/lspaccatrosi16/tmappr/lib/types"
	"github.com/lspaccatrosi16/tmappr/lib/util"
)

func doOverview(counter *int, iptLines *[]string) (*map[string]*types.Line, []*types.Line, error) {
	logger := logging.GetLogger()

	if err := expectLine((*iptLines)[*counter], "[Lines]"); err != nil {
		return nil, nil, err
	}

	(*counter)++

	util.DebugSection("Parsing Lines Overview")

	lines := []*types.Line{}

	for {
		if *counter >= len(*iptLines) {
			break
		}

		if (*iptLines)[*counter] == "" {
			(*counter)++
			continue
		}

		if strings.HasPrefix((*iptLines)[*counter], "[") {
			break
		}

		parsed, err := types.ParseLine((*iptLines)[*counter])
		if err != nil {
			return nil, nil, err
		}

		logger.Debug(fmt.Sprintf("Parsed Line Overview %s", parsed.Code))

		lines = append(lines, parsed)
		(*counter)++
	}

	lineMap := map[string]*types.Line{}

	for _, l := range lines {
		lineMap[l.Code] = l
	}

	return &lineMap, lines, nil
}
