package data

import (
	"errors"
	"fmt"
	"net/http"

	"{{.ImportPath}}/helpers"
)

var (
	ErrFailedValidation = errors.New("Failed validation")
	ErrFailedDecodingInput = errors.New("Failed decoding input")
	ErrFailedValidatingInput = errors.New("Failed validating input")
	ErrFailedConvertingInput = errors.New("Failed converting input")
	ErrFailedTransformingInput = errors.New("Failed transforming input")
)

{{range $message := .Messages}}
type {{$message.TitleCamel}} struct {
	{{range $field := $message.Fields}}{{$field.TitleCamel}} {{$field.DataType}} `json:"{{$field.Snake}},omitempty" validate:"{{$field.Validate}}" transform:"{{$field.Transform}}"`
	{{end}}
}
{{end}}

{{range $message := .Messages}}{{if $message.IsInput}}
func Get{{$message.TitleCamel}}(r *http.Request) ({{$message.Camel}} {{$message.TitleCamel}}, err error) {
	inputP := &{{$message.TitleCamel}}P{}
	if err := helpers.Unmarshal(r, inputP); err != nil {
		fmt.Println("Failed decoding input:", err)
		return {{$message.Camel}}, ErrFailedDecodingInput
	}

	msg, err := helpers.Validate(inputP)
	if err != nil {
		fmt.Println("Failed validating input:", err)
		return {{$message.Camel}}, ErrFailedValidatingInput
	}

	if msg != "" {
		fmt.Println("Failed validation:", msg)
		return {{$message.Camel}}, ErrFailedValidation
	}	

	if err := helpers.Transform(inputP); err != nil {
		fmt.Println("Failed transforming input:", err)
		return {{$message.Camel}}, ErrFailedTransformingInput
	}	

	{{$message.Camel}} = Convert{{$message.TitleCamel}}P(inputP)

	return
}
{{end}}{{end}}


