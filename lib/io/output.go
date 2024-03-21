package io

import (
	"bytes"
	"io"
	"os"

	"github.com/lspaccatrosi16/tmappr/lib/types"
)

func OutputFile(config *types.AppConfig, data *bytes.Buffer) error {
	f, err := os.Create(config.Output)
	if err != nil {
		return err
	}
	defer f.Close()

	io.Copy(f, data)
	return nil
}
