package io

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/lspaccatrosi16/go-cli-tools/logging"
	"github.com/lspaccatrosi16/tmappr/lib/types"
)

func ParseFile(config *types.AppConfig) (map[string]*types.Line, []*types.Stop, error) {
	logger := logging.GetLogger()

	var endStr string
	if config.Ending == types.CRLF {
		endStr = "\r\n"
	} else {
		endStr = "\n"
	}

	f, err := os.Open(config.Input)

	if err != nil {
		return nil, nil, err
	}

	defer f.Close()

	buf := bytes.NewBuffer(nil)
	io.Copy(buf, f)

	iptLines := strings.Split(buf.String(), endStr)
	if expectLine(iptLines[0], "[Lines]") != nil {
		return nil, nil, err
	}

	counter := 1

	lines := []*types.Line{}
	stops := []*types.Stop{}

	logger.Debug("Parsing Lines")

	for {
		if counter >= len(iptLines) {
			return nil, nil, fmt.Errorf("unexpected EOF")
		}

		if iptLines[counter] == "" {
			counter++
			continue
		}

		if expectLine(iptLines[counter], "[Stops]") == nil {
			counter++
			break
		}

		parsed, err := types.ParseLine(iptLines[counter])
		if err != nil {
			return nil, nil, err
		}

		lines = append(lines, parsed)
		counter++
	}

	logger.Debug("Parsing Stops")

	for {
		if counter >= len(iptLines) {
			break
		}

		if iptLines[counter] == "" {
			counter++
			continue
		}

		parsed, err := types.ParseStop(iptLines[counter])
		if err != nil {
			return nil, nil, err
		}

		stops = append(stops, parsed)
		counter++
	}

	lineMap := map[string]*types.Line{}

	for _, l := range lines {
		lineMap[l.Code] = l
	}

	return lineMap, stops, nil
}

func expectLine(line, str string) error {
	if strings.HasPrefix(line, str) {
		return nil
	}
	return fmt.Errorf("expected %s", str)
}
