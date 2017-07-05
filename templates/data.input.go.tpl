package data

import (
	"errors"
	"net/http"

	"{{.ImportPath}}/helpers"
	log "github.com/sirupsen/logrus"
)

var (
	ErrFailedDecodingInput = errors.New("Failed decoding input")
)

{{range $input := .Inputs}}
type {{$input.TitleCamel}} struct {
	{{range $field := $input.Fields}}{{$field.TitleCamel}} {{$field.DataType}} `json:"{{$field.Snake}},omitempty" validate:"{{$field.Validate}}" transform:"{{$field.Transform}}"`
	{{end}}
}
{{end}}

{{range $message := .Inputs}}{{if $message.RPCInput}}
func Get{{MinusP $message.TitleCamel}}(r *http.Request) ({{MinusP $message.Camel}} {{MinusP $message.TitleCamel}}, err error) {
	inputP := &{{$message.TitleCamel}}{}
	if err := helpers.Unmarshal(r, inputP); err != nil {
		log.Error("Failed decoding input:", err)
		return {{MinusP $message.Camel}}, ErrFailedDecodingInput
	}

	msg, err := helpers.Validate(inputP)
	if err != nil {
		return {{MinusP $message.Camel}}, err
	}

	if msg != "" {
		return {{MinusP $message.Camel}}, errors.New(msg)
	}	

	if err := helpers.Transform(inputP); err != nil {
		return {{MinusP $message.Camel}}, err
	}	

	{{MinusP $message.Camel}} = Convert{{$message.TitleCamel}}(inputP)

	return
}
{{end}}{{end}}

{{range $input := .Inputs}}
func Convert{{$input.TitleCamel}}({{$input.Camel}} *{{$input.TitleCamel}}) ({{MinusP $input.Camel}} {{MinusP $input.TitleCamel}}) {
	{{range $field := $input.Fields}}
	{{if $field.IsRepeatedStruct}}

	{{$field.Camel}} := {{MinusStar $field.DataType | MinusP}}{}
	for _, field := range {{$input.Camel}}.{{$field.TitleCamel}} {
		{{$field.Camel}} = append({{$field.Camel}}, Convert{{$field.DataTypeName.TitleCamel}}P(field))
	}
	{{MinusP $input.Camel}}.{{$field.TitleCamel}} = {{$field.Camel}}

	{{else if $field.IsStruct}}

	if {{$input.Camel}}.{{$field.TitleCamel}} != nil {
		{{MinusP $input.Camel}}.{{$field.TitleCamel}} = Convert{{$field.DataTypeName.TitleCamel}}P({{$input.Camel}}.{{$field.TitleCamel}})
	}

	{{else if $field.IsRepeated}}

	{{$field.Camel}} := {{MinusStar $field.DataType}}{}
	for _, field := range {{$input.Camel}}.{{$field.TitleCamel}} {
		{{$field.Camel}} = append({{$field.Camel}}, *field)
	}
	{{MinusP $input.Camel}}.{{$field.TitleCamel}} = {{$field.Camel}}


	{{else}}

	if {{$input.Camel}}.{{$field.TitleCamel}} != nil {
		{{MinusP $input.Camel}}.{{$field.TitleCamel}} = *{{$input.Camel}}.{{$field.TitleCamel}}
	}

	{{end}}{{end}}
	return 
}
{{end}}
