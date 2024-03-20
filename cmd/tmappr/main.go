package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/lspaccatrosi16/go-cli-tools/logging"
	"github.com/lspaccatrosi16/tmappr/lib/engine"
	"github.com/lspaccatrosi16/tmappr/lib/flags"
	"github.com/lspaccatrosi16/tmappr/lib/io"
)

func main() {
	logger := logging.GetLogger()

	config, err := flags.GetFlagData()
	handle(err)

	logger.SetVerbose(config.Verbose)
	logger.Debug("Parsing Input File")

	lineMap, stopMap, err := io.ParseFile(config)
	handle(err)

	lineErrors := []string{}

	var maxX, maxY float64

	for _, s := range *stopMap {
		if s.Coordinates[0] > maxX {
			maxX = s.Coordinates[0]
		}

		if s.Coordinates[1] > maxY {
			maxY = s.Coordinates[1]
		}
	}

	if len(lineErrors) > 0 {
		handle(fmt.Errorf("could not find line(s) with code %s", strings.Join(lineErrors, ", ")))
	}

	err = engine.RunEngine(config, lineMap, stopMap, [2]float64{float64(config.XRes) / maxY, float64(config.YRes) / maxY})
	handle(err)

}

func handle(e error) {
	if e != nil {
		fmt.Println(e.Error())
		os.Exit(1)
	}
}
