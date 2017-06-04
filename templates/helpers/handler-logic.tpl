{{$.Input.Camel}}, err := data.Get{{$.Input.TitleCamel}}(r)
if err != nil {
    return
}
fmt.Println({{$.Input.Camel}})


// {{$.Output.Camel}} := data.{{$.Output.TitleCamel}}{}

// jsonBytes, err := json.Marshal(output)
// if err != nil {
//     fmt.Println("Failed marshalling to JSON:", err)
//     http.Error(w, ErrorJSON("JSON Marshal Error"), http.StatusInternalServerError)
//     return
// }
// 
// if _, err := w.Write(jsonBytes); err != nil {
//     fmt.Println("Failed writing to response writer:", err)
//     http.Error(w, ErrorJSON("Failed writing to output"), http.StatusInternalServerError)
//     return
// }
