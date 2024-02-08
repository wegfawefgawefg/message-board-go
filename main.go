package main

import (
	"html/template"
	"log"
	"net/http"
)

const saved_pages_path = "pages"
const templates_path = "templates"

type Page struct {
	Title string
	Body  []byte
}

var templates = template.Must(
	template.ParseFiles(
		templates_path+"/edit.html",
		templates_path+"/view.html",
		templates_path+"/index.html",
		templates_path+"/add_new_page.html",
	))

func main() {
	http.HandleFunc("/view/", makeHandler(viewHandler))
	http.HandleFunc("/edit/", makeHandler(editHandler))
	http.HandleFunc("/save/", makeHandler(saveHandler))
	http.HandleFunc("/add_new_page", addNewPageHandler)
	http.HandleFunc("/internal_add_new_page", internalAddNewPageHandler)
	http.HandleFunc("/index", indexHandler)
	http.HandleFunc("/", rootHandler)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
