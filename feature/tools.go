package feature

import "regexp"

var leadZeroRegext = regexp.MustCompile("(:)(0+)(\\d+)")

func ClearLeadZero(str string) string {
	return leadZeroRegext.ReplaceAllString(str, "$1$3")
}