package data

{{range $input := .Inputs}}
type {{$input.TitleCamel}} struct {
	{{range $field := $input.Fields}}{{$field.TitleCamel}} {{$field.DataType}} `json:"{{$field.Snake}},omitempty" validate:"{{$field.Validate}}" transform:"{{$field.Transform}}"`
	{{end}}
}
{{end}}
