package generate

import (
	"strings"

	log "github.com/sirupsen/logrus"
)

const (
	TransformStrEncrypt      = "encrypt"
	TransformStrDecrypt      = "decrypt"
	TransformStrHash         = "hash"
	TransformStrPasswordHash = "passwordHash"
	TransformStrTruncate     = "truncate"
	TransformStrTrimChars    = "trimChars"
	TransformStrTrimSpace    = "trimSpace"
	TransformStrDefault      = "default"
	ValidateStrMaxLength     = "maxLength"
	ValidateStrMinLength     = "minLength"
	ValidateStrGreaterThan   = "greaterThan"
	ValidateStrLessThan      = "lessThan"
	ValidateStrRequired      = "required"
	ValidateStrMustHaveChars = "mustHaveChars"
	ValidateStrCantHaveChars = "cantHaveChars"
	ValidateStrOnlyHaveChars = "onlyHaveChars"
	DataTypeFloat32          = "float32"
	DataTypeFloat64          = "float64"
	DataTypeInt              = "int"
	DataTypeInt32            = "int32"
	DataTypeInt64            = "int64"
	DataTypeString           = "string"
	DataTypeBool             = "bool"
	DataTypeStringArr        = "[]string"
	DataTypeIntArr           = "[]int"
	DataTypeInt32Arr         = "[]int32"
	DataTypeInt64Arr         = "[]int64"
	DataTypeFloat32Arr       = "[]float32"
	DataTypeFloat64Arr       = "[]float64"
	DataTypeBoolArr          = "[]bool"
)

type Template struct {
	TemplateName string
	FileName     string
	Data         interface{}
}

func HandleQuotes(value, typeStr string) string {
	if strings.ToLower(typeStr) == DataTypeString {
		return `"` + value + `"`
	}
	return value
}

func EmptyValue(dataType string) (out string) {
	if dataType[:2] == "[]" {
		return dataType + "{}"
	}

	switch dataType {
	case DataTypeString:
		return "\"\""
	case DataTypeInt:
		fallthrough
	case DataTypeInt32:
		fallthrough
	case DataTypeInt64:
		return "0"
	case DataTypeFloat32:
		fallthrough
	case DataTypeFloat64:
		return "0.0"
	case DataTypeBool:
		return "false"
	}
	log.Warn("DATA TYPE NOT DEFINED: ", dataType)
	return dataType + "{}"
}

func GetHTTPMethod(method string) (httpMethod string) {
	method = strings.ToLower(method)
	switch method {
	case "connect":
		return "MethodConnect"
	case "delete":
		return "MethodDelete"
	case "get":
		return "MethodGet"
	case "head":
		return "MethodHead"
	case "options":
		return "MethodOptions"
	case "patch":
		return "MethodPatch"
	case "post":
		return "MethodPost"
	case "put":
		return "MethodPut"
	case "trace":
		return "MethodTrace"
	}
	log.Warn("BAD METHOD PROVIDED: ", method)
	return
}

func FallbackSet(in, fallback string) string {
	if in != "" {
		return in
	}
	return fallback
}

func GetMethodMiddlewares(in string, cfg Degeneres) (out string) {
	// mws := []string{}
	// for _, path := range cfg.API.Paths {
	// 	for _, method := range path.Methods {
	// 		if method.Name == in {
	// 			for k := range method.Middlewares {
	// 				mws = append(mws, k)
	// 			}
	// 			break
	// 		}
	// 	}
	// }

	// out = ""
	// for _, mw := range mws {
	// 	out += ", Middleware" + mw
	// }

	return
}

func GetPathMiddlewares(cfg Degeneres) (out string) {
	// mws := []string{}
	// for k := range cfg.API.Middlewares {
	// 	mws = append(mws, k)
	// }

	// out = ""
	// for _, mw := range mws {
	// 	out += ", Middleware" + mw
	// }
	return
}

func GetInputType(inputType string) (out string) {
	if len(inputType) < 2 {
		return inputType
	}
	if inputType[:2] != "[]" {
		return "*" + inputType
	}
	return "[]*" + inputType[2:]
}

func GetDereferenceFunc(outputType string) (out string) {
	switch outputType {
	case DataTypeStringArr:
		return "dereferenceStringArray"
	case DataTypeIntArr:
		return "dereferenceIntArray"
	case DataTypeInt32Arr:
		return "dereferenceInt32Array"
	case DataTypeInt64Arr:
		return "dereferenceInt64Array"
	case DataTypeFloat32Arr:
		return "dereferenceFloat32Array"
	case DataTypeFloat64Arr:
		return "dereferenceFloat64Array"
	case DataTypeBoolArr:
		return "dereferenceBoolArray"
	}

	return "*"
}

func IsStruct(dataType string) (isStruct bool) {
	return !(dataType == DataTypeFloat32 ||
		dataType == DataTypeFloat64 ||
		dataType == DataTypeInt ||
		dataType == DataTypeInt32 ||
		dataType == DataTypeInt64 ||
		dataType == DataTypeString ||
		dataType == DataTypeBool ||
		dataType == DataTypeStringArr ||
		dataType == DataTypeIntArr ||
		dataType == DataTypeInt32Arr ||
		dataType == DataTypeInt64Arr ||
		dataType == DataTypeFloat32Arr ||
		dataType == DataTypeFloat64Arr ||
		dataType == DataTypeBoolArr)
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
