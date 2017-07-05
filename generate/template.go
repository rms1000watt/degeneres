package generate

import (
	"strings"
)

type Template struct {
	TemplateName string
	FileName     string
	Data         interface{}
}

func MinusP(in string) string {
	if string(in[len(in)-1]) == "P" {
		return in[:len(in)-1]
	}
	return in
}

func MinusStar(in string) string {
	return strings.Replace(in, "*", "", -1)
}
