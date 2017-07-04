{{$.Input.Camel}}, err := data.Get{{$.Input.TitleCamel}}(r)
if err != nil {
    log.Error(err)
    http.Error(w, helpers.ErrorJSON(err.Error()), http.StatusInternalServerError)
    return
}
log.Debug({{$.Input.Camel}})

{{$.Output.Camel}} := data.{{$.Output.TitleCamel}}{}

// Developer do stuff here

helpers.WriteJSON({{$.Output.Camel}}, w)
