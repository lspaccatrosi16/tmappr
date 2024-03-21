package main

import (
	"fmt"
	"os"

	"github.com/lspaccatrosi16/go-cli-tools/logging"
	"github.com/lspaccatrosi16/tmappr/lib/drawer"
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

	pathings, cStopMap, combinedGrid, maxX, maxY, err := engine.RunEngine(config, lineMap, stopMap)
	handle(err)

	drawn := drawer.DrawMap(config, pathings, &cStopMap, lineMap, combinedGrid, maxX, maxY)

	err = io.OutputFile(config, drawn)
	handle(err)

}

func handle(e error) {
	if e != nil {
		fmt.Println(e.Error())
		os.Exit(1)
	}
}
