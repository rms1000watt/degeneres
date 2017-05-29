{{range $option := .API.Middlewares.CORS.Options}}
var (
    {{if eq $option.Key "Hosts"}}CORSHosts = []string{ {{$option.Value}} }
    {{end}}
)
{{end}}