addr := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)

fmt.Println("Starting server at:", addr)
{{if .CertsPath}}log.Fatal(http.ListenAndServeTLS(addr, filepath.Join(cfg.CertsPath, cfg.CertName), filepath.Join(cfg.CertsPath, cfg.KeyName), ServerHandler()))
{{else}}log.Fatal(http.ListenAndServe(addr, ServerHandler()))
{{end}}