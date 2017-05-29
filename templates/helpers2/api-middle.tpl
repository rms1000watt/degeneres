{{if .CommandLine.Command.API}}addr := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)

fmt.Println("Starting server at:", addr)
{{if .API.CertsPath}}log.Fatal(http.ListenAndServeTLS(addr, "./certs/{{FallbackSet .API.PubKeyFileName "server.cer"}}", "./certs/{{FallbackSet .API.PrivKeyFileName "server.key"}}", ServerHandler()))
{{else}}log.Fatal(http.ListenAndServe(addr, ServerHandler()))
{{end}}
{{end}}