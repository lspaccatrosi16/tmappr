package util

import "github.com/lspaccatrosi16/go-cli-tools/logging"

func DebugSection(s string) {
	logger := logging.GetLogger()

	logger.Debug("")
	logger.Debug(s)
	logger.DebugDivider()
}
