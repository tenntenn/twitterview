package main

import (
    "net/http"
    "html/template"
)

var templates = template.Must(template.ParseGlob("template/*.html"))

func main() {


    index := func(w http.ResponseWriter, r *http.Request) {
        templates.ExecuteTemplate(w, "index", nil)
    }

    mux := http.NewServeMux()
    mux.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("css"))))
    mux.Handle("/js/", http.StripPrefix("/js/", http.FileServer(http.Dir("js"))))
    mux.Handle("/lib/", http.StripPrefix("/lib/", http.FileServer(http.Dir("lib"))))
    mux.HandleFunc("/", index)
    http.ListenAndServe(":9000", mux)
    //l, _ := net.Listen("tcp", ":9000")
    //fcgi.Serve(l, mux)
}
