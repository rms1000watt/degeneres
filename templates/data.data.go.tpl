package data

import (
	"database/sql"

	"github.com/jmoiron/sqlx/types"
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
