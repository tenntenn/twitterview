package main

import (
    "net/http"
    "html/template"
    "strings"
    "time"
    "math/rand"
    "fmt"
    "io"
    "os"
)

var templates = template.Must(template.ParseGlob("template/*.html"))


func index(w http.ResponseWriter, r *http.Request) {
    templates.ExecuteTemplate(w, "index", nil)
}

var random = rand.New(rand.NewSource(time.Now().UnixNano()))

func create(w http.ResponseWriter, r *http.Request) {

    // Form data
    login := r.FormValue("login")
    password := r.FormValue("password")
    keywords := strings.Split(strings.Replace(r.FormValue("keywords"), " ", "", -1), ",")
    file, _, _ := r.FormFile("pdf")

    // create view
    seed := fmt.Sprintf("%s %s %v %d", login, password, keywords, random.Int())
    view := NewView(seed)

    // save pdf file
    path := fmt.Sprintf("pdf/%s.pdf", view.Id)
    pdfFile, _ := os.Create(path)
    defer pdfFile.Close()
    io.Copy(pdfFile, file)

    // redirect view
    view.Start()
    viewPath := fmt.Sprintf("/view/%s/show", view.Id)
    http.Redirect(w, r, viewPath, http.StatusFound)
}

var mux = http.NewServeMux()

func main() {
    mux.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("css"))))
    mux.Handle("/js/", http.StripPrefix("/js/", http.FileServer(http.Dir("js"))))
    mux.Handle("/lib/", http.StripPrefix("/lib/", http.FileServer(http.Dir("lib"))))
    mux.Handle("/pdf/", http.StripPrefix("/pdf/", http.FileServer(http.Dir("pdf"))))
    mux.HandleFunc("/create", create)
    mux.HandleFunc("/", index)
    http.ListenAndServe(":9000", mux)
    //l, _ := net.Listen("tcp", ":9000")
    //fcgi.Serve(l, mux)
}
