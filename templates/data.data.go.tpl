package data

import (
	"database/sql"

	"github.com/jmoiron/sqlx/types"
	log "github.com/sirupsen/logrus"
)

{{range $message := .Messages}}
type {{$message.TitleCamel}} struct {
	{{- range $field := $message.Fields}}
	{{$field.TitleCamel}} {{$field.DataType}} `json:"{{$field.Snake}},omitempty"`
	{{- end}}
}
{{end}}

{{range $message := .Messages}}
type {{$message.TitleCamel}}DB struct {
	{{- range $field := $message.Fields}}
	{{$field.TitleCamel}} {{$field.DataTypeDB}} `db:"{{$field.Snake}}"`
	{{- end}}
}
{{end}}

{{range $message := .Messages}}
func Convert{{AddDB $message.TitleCamel}}({{AddDB $message.Camel}} {{AddDB $message.TitleCamel}}) ({{$message.Camel}} {{$message.TitleCamel}}) {
	{{range $field := $message.Fields}}
		{{if $field.IsRepeatedStruct}}
			{{$field.Camel}} := {{$field.DataType}}{}
			for _, field := range {{AddDB $message.Camel}}.{{$field.TitleCamel}} {
				{{$field.Camel}} = append({{$field.Camel}}, Convert{{AddDB $field.DataTypeName.TitleCamel}}(field))
			}
			{{$message.Camel}}.{{$field.TitleCamel}} = {{$field.Camel}}

		{{else if $field.IsStruct}}
			{{$message.Camel}}.{{$field.TitleCamel}} = Convert{{AddDB $field.DataTypeName.TitleCamel}}({{AddDB $message.Camel}}.{{$field.TitleCamel}})

		{{else if $field.IsRepeated}}
			{{$field.Camel}} := {{$field.DataType}}{}
			for _, field := range {{$message.Camel}}.{{$field.TitleCamel}} {
				{{$field.Camel}} = append({{$field.Camel}}, field)
			}
			{{$message.Camel}}.{{$field.TitleCamel}} = {{$field.Camel}}

		{{else}}
			{{if IsMap $field.DataType}}
				{{$field.Camel}} := {{$field.DataType}}{}
				if err := {{AddDB $message.Camel}}.{{$field.TitleCamel}}.Unmarshal(&{{$field.Camel}}); err != nil {
					log.Error("Failed unmarshalling {{$field.Camel}}:", err)
				} else {
					{{$message.Camel}}.{{$field.TitleCamel}} = {{$field.Camel}}
				}

			{{else}}
				if {{AddDB $message.Camel}}.{{$field.TitleCamel}}.Valid {
					{{$message.Camel}}.{{$field.TitleCamel}} = {{AddDB $message.Camel}}.{{$field.TitleCamel}}.{{ConvertFromDBDataType $field.DataType}}
				}

			{{end}}
		{{end}}
	{{end}}
	return
}
{{end}}
