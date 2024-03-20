package main

import (
	"fmt"
	"os"

	"github.com/lspaccatrosi16/go-cli-tools/logging"
	"github.com/lspaccatrosi16/tmappr/lib/flags"
	"github.com/lspaccatrosi16/tmappr/lib/io"
)

func main() {
	logger := logging.GetLogger()

	config, err := flags.GetFlagData()
	handle(err)

	logger.SetVerbose(config.Verbose)
	logger.Debug("Parsing Input File")

	lines, stops, err := io.ParseFile(config)
	handle(err)

	fmt.Println(lines)
	fmt.Println(stops)
}

func handle(e error) {
	if e != nil {
		fmt.Println(e.Error())
		os.Exit(1)
	}
}
