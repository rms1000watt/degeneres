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

func AddDB(in string) string {
	return in + "DB"
}

func ConvertFromDBDataType(in string) string {
	if isDbDataType(in) {
		if isInt(in) {
			return "Int64"
		}
		if isFloat(in) {
			return "Float64"
		}
		if in == "string" {
			return "String"
		}
		if in == "bool" {
			return "Bool"
		}
	}

	return in
}

func IsMap(in string) bool {
	return strings.Contains(in, "map[")
}

func MinusStar(in string) string {
	return strings.Replace(in, "*", "", -1)
}
