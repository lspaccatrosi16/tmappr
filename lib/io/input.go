package io

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/lspaccatrosi16/tmappr/lib/types"
)

func ParseFile(config *types.AppConfig, data *types.AppData) error {
	var endStr string
	if config.Ending == types.CRLF {
		endStr = "\r\n"
	} else {
		endStr = "\n"
	}

	f, err := os.Open(config.Input)

	if err != nil {
		return err
	}

	defer f.Close()

	buf := bytes.NewBuffer(nil)
	io.Copy(buf, f)

	iptLines := strings.Split(buf.String(), endStr)
	counter := 0

	lineMap, lines, err := doOverview(&counter, &iptLines)
	if err != nil {
		return err
	}

	stopMap, stops, err := doStops(&counter, &iptLines)
	if err != nil {
		return err
	}

	err = doLines(&counter, &iptLines, lineMap, stopMap)
	if err != nil {
		return err
	}

	stopLineMap := map[*types.Stop][]*types.Line{}

	for _, line := range *lineMap {
		for _, stop := range line.Stops {
			if _, ok := stopLineMap[stop]; !ok {
				stopLineMap[stop] = []*types.Line{}
			}
			stopLineMap[stop] = append(stopLineMap[stop], line)
		}
	}

	data.LinesNames = *lineMap
	data.StopNames = *stopMap
	data.Stops = stops
	data.Lines = lines
	data.StopLineMap = stopLineMap

	return nil
}

func expectLine(line, str string) error {
	if strings.HasPrefix(line, str) {
		return nil
	}
	return fmt.Errorf("expected\n%s\nbut got\n%s", str, line)
}
