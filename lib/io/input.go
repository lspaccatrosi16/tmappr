package io

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/lspaccatrosi16/tmappr/lib/types"
)

func ParseFile(config *types.AppConfig) (*map[string]*types.Line, *map[string]*types.Stop, error) {
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
	counter := 0

	lineMap, err := doOverview(&counter, &iptLines)
	if err != nil {
		return nil, nil, err
	}

	stopMap, err := doStops(&counter, &iptLines)
	if err != nil {
		return nil, nil, err
	}

	err = doLines(&counter, &iptLines, lineMap, stopMap)
	if err != nil {
		return nil, nil, err
	}

	return lineMap, stopMap, nil
}

func expectLine(line, str string) error {
	if strings.HasPrefix(line, str) {
		return nil
	}
	return fmt.Errorf("expected\n%s\nbut got\n%s", str, line)
}
