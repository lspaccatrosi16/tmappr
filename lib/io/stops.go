package io

import (
	"fmt"
	"strings"

	"github.com/lspaccatrosi16/go-cli-tools/logging"
	"github.com/lspaccatrosi16/tmappr/lib/types"
	"github.com/lspaccatrosi16/tmappr/lib/util"
)

func doStops(counter *int, iptLines *[]string) (*map[string]*types.Stop, error) {
	logger := logging.GetLogger()

	if err := expectLine((*iptLines)[*counter], "[Stops]"); err != nil {
		return nil, err
	}

	(*counter)++

	util.DebugSection("Parsing Stops")

	stops := []*types.Stop{}

	stopId := 1

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

		parsed, err := types.ParseStop((*iptLines)[*counter])
		if err != nil {
			return nil, err
		}

		parsed.Id = stopId

		logger.Debug(fmt.Sprintf("Parsed Stop %s ID: %d", parsed.Name, parsed.Id))

		stops = append(stops, parsed)
		(*counter)++
		stopId++
	}

	stopMap := map[string]*types.Stop{}

	for _, s := range stops {
		stopMap[s.Name] = s
	}

	return &stopMap, nil
}
