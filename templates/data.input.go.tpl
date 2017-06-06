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

{{range $input := .Inputs}}
type {{$input.TitleCamel}} struct {
	{{range $field := $input.Fields}}{{$field.TitleCamel}} {{$field.DataType}} `json:"{{$field.Snake}},omitempty" validate:"{{$field.Validate}}" transform:"{{$field.Transform}}"`
	{{end}}
}
{{end}}

{{range $message := .Inputs}}
func Get{{MinusP $message.TitleCamel}}(r *http.Request) ({{$message.Camel}} {{$message.TitleCamel}}, err error) {
	inputP := &{{$message.TitleCamel}}{}
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

	{{$message.Camel}} = Convert{{$message.TitleCamel}}(inputP)

	return
}
{{end}}

{{range $input := .Inputs}}
func Convert{{$input.TitleCamel}}({{$input.Camel}} *{{$input.TitleCamel}}) ({{$input.Camel}} {{$input.TitleCamel}}) {
	{{range $field := $input.Fields}}{{if $field.IsStruct}}

	if {{$input.Camel}}.{{$field.TitleCamel}} != nil {
		{{$input.Camel}}.{{$field.TitleCamel}} = Convert{{$field.TitleCamel}}P({{$input.Camel}}.{{$field.TitleCamel}})
	}

	{{else if $field.IsRepeated}}

	{{$field.Camel}} := {{$field.DataType}}{}
	for _, field := range {{$input.Camel}}.{{$field.TitleCamel}} {
		{{$field.Camel}} = append({{$field.Camel}}, *field)
	}
	{{$input.Camel}}.{{$field.TitleCamel}} = {{$field.Camel}}

	{{else if $field.IsRepeatedStruct}}



	{{else}}

	if {{$input.Camel}}.{{$field.TitleCamel}} != nil {
		{{$input.Camel}}.{{$field.TitleCamel}} = *{{$input.Camel}}.{{$field.TitleCamel}}
	}

	{{end}}{{end}}
	return 
}
{{end}}
