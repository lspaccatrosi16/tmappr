package util

import "strings"

func Trim(s string) string {
	return strings.Trim(s, " \t\r\n")
}
