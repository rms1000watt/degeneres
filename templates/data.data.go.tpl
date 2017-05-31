package data

{{range $message := .Messages}}
type {{$message.TitleCamel}} struct {
	{{range $field := $message.Fields}}{{$field.TitleCamel}} {{$field.DataType}} `json:"{{$field.Snake}},omitempty" validate:"{{$field.Validate}}" transform:"{{$field.Transform}}"`
	{{end}}
}
{{end}}

{{`/*
{{range $path := .API.Paths}}
{{range $method := $path.Methods}}
func get{{$path.Name | Title}}Output{{$method.Name | ToUpper}}({{$path.Name | ToLower}}Input{{$method.Name | ToUpper}} *{{$path.Name | Title}}Input{{$method.Name | ToUpper}}) ({{$path.Name | ToLower}}Output{{$method.Name | ToUpper}} {{$path.Name | Title}}Output{{$method.Name | ToUpper}}) {
	if {{$path.Name | ToLower}}Input{{$method.Name | ToUpper}} == nil {
		return
	}
	
	{{range $output := $method.Outputs}}{{if OutputInInputs $output.Name $method.Inputs}}{{$output.Name | ToCamelCase}} := {{EmptyValue $output.Type}}
	if {{$path.Name | ToLower}}Input{{$method.Name | ToUpper}}.{{$output.Name | Title}} != nil {
		{{$output.Name | ToCamelCase}} = {{GetDereferenceFunc $output.Type}}({{$path.Name | ToLower}}Input{{$method.Name | ToUpper}}.{{$output.Name | Title}})
	}{{end}}
	
	{{end}}

	{{$path.Name | ToLower}}Output{{$method.Name | ToUpper}} = {{$path.Name | Title}}Output{{$method.Name | ToUpper}}{
		{{range $output := $method.Outputs}}{{if OutputInInputs $output.Name $method.Inputs}}{{$output.Name | Title}}: {{$output.Name | ToCamelCase}},
		{{end}}{{end}}
	}
	return
}
{{end}}
{{end}}
*/`}}
