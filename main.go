package main

import (
	"html/template"
	"net/http"
)

var tpl *template.Template

type PageData struct {
	ErrorMessage string
}

func init() {
	tpl = template.Must(template.ParseGlob("templates/*.gohtml"))
}

func main() {
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/new", newHandler)
	http.ListenAndServe(":8080", nil)
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	data := PageData{}

	data.ErrorMessage = "No error"

	tpl.ExecuteTemplate(w, "index.gohtml", data)
}

func newHandler(w http.ResponseWriter, r *http.Request) {
	data := PageData{}

	data.ErrorMessage = "No error"

	tpl.ExecuteTemplate(w, "new.gohtml", data)
}
