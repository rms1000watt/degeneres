package {{.CommandLine.Command.Name}}

{{range $path := .API.Paths}}
{{range $method := $path.Methods}}
type {{$path.Name | Title}}Input{{$method.Name | ToUpper}} struct {
    {{range $input := $method.Inputs}}{{$input.Name}} {{GetInputType $input.Type}} `json:"{{$input.DisplayName}},omitempty" validate:"{{GenValidationStr $input}}" transform:"{{GenTransformStr $input}}"`
    {{end}}
}

type {{$path.Name | Title}}Output{{$method.Name | ToUpper}} struct {
    {{range $output := $method.Outputs}}{{$output.Name}} {{$output.Type}} `json:"{{$output.DisplayName}},omitempty"`
    {{end}}
}

{{GetStructs2 $method.Inputs $method.Outputs $.Structs}}
{{end}}
{{end}}

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
