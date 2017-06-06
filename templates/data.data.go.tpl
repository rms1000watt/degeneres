package data

{{range $message := .Messages}}
type {{$message.TitleCamel}} struct {
	{{range $field := $message.Fields}}{{$field.TitleCamel}} {{$field.DataType}} `json:"{{$field.Snake}},omitempty"`
	{{end}}
}
{{end}}
