package main

import (
	"fmt"
	"os"

	"github.com/lspaccatrosi16/go-cli-tools/logging"
	"github.com/lspaccatrosi16/tmappr/lib/drawer"
	"github.com/lspaccatrosi16/tmappr/lib/engine"
	"github.com/lspaccatrosi16/tmappr/lib/flags"
	"github.com/lspaccatrosi16/tmappr/lib/io"
	"github.com/lspaccatrosi16/tmappr/lib/types"
)

func main() {
	logger := logging.GetLogger()

	config, err := flags.GetFlagData()
	handle(err)

	data := types.AppData{}

	logger.SetVerbose(config.Verbose)
	logger.Debug("Parsing Input File")

	err = io.ParseFile(config, &data)
	handle(err)

	err = engine.RunEngine(config, &data)
	handle(err)

	drawn := drawer.DrawMap(config, &data)

	err = io.OutputFile(config, drawn)
	handle(err)
}

func handle(e error) {
	if e != nil {
		fmt.Println(e.Error())
		os.Exit(1)
	}
}
