package data

{{range $input := .Inputs}}
type {{$input.TitleCamel}} struct {
	{{range $field := $input.Fields}}{{$field.TitleCamel}} {{$field.DataType}} `json:"{{$field.Snake}},omitempty" validate:"{{$field.Validate}}" transform:"{{$field.Transform}}"`
	{{end}}
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
