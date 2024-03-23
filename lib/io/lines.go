package io

import (
	"fmt"
	"strings"

	"github.com/lspaccatrosi16/go-cli-tools/logging"
	"github.com/lspaccatrosi16/tmappr/lib/types"
	"github.com/lspaccatrosi16/tmappr/lib/util"
)

func doLines(counter *int, iptLines *[]string, lineMap *map[string]*types.Line, stopMap *map[string]*types.Stop) error {
	logger := logging.GetLogger()

	util.DebugSection("Parsing Individual Lines")

	for {
		if *counter >= len(*iptLines) {
			break
		}

		codeLine := (*iptLines)[*counter]
		if util.Trim(codeLine) == "" {
			(*counter)++
			continue
		}

		if !strings.HasPrefix(codeLine, "[") {
			break
		}

		if !strings.HasSuffix(codeLine, "]") {
			return fmt.Errorf("expected a line decleration surrounded by brackets")
		}

		lineCode := codeLine[1 : len(codeLine)-1]
		(*counter)++

		stopsLine := (*iptLines)[*counter]
		stops := strings.Split(stopsLine, ",")

		lineStops := []*types.Stop{}

		for i, s := range stops {
			stop := util.Trim(s)
			if !strings.HasPrefix(stop, "\"") || !strings.HasSuffix(stop, "\"") {
				return fmt.Errorf("entry name must be double quoted (%s)", s)
			}

			stopName := stop[1 : len(stop)-1]

			s, ok := (*stopMap)[stopName]

			if !ok {
				return fmt.Errorf("unknown stop, \"%s\"", stopName)
			}

			s.Lines = append(s.Lines, lineCode)

			if i == 0 || i == len(stops)-1 {
				s.Type = types.Interchange
			}

			lineStops = append(lineStops, s)
		}

		lEntry := (*lineMap)[lineCode]
		lEntry.Stops = lineStops

		(*counter)++

		logger.Debug(fmt.Sprintf("Parsed Line %s", lineCode))
	}

	return nil
}
