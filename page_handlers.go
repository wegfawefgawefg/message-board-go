package main

import (
	"fmt"
	"net/http"
	"os"
	"regexp"
)

// //////////////////////	Handlers ///////////////////////

func rootHandler(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/index", http.StatusFound)
}

// struct with two strings, display title and title
type IndexPageListing struct {
	DisplayTitle string
	Title        string
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("pageIndexHandler start")
	files, err := os.ReadDir(saved_pages_path)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var page_listings []IndexPageListing
	for _, file := range files {
		page_listing := IndexPageListing{
			DisplayTitle: filenameToTitle(file.Name()),
			Title:        file.Name(),
		}
		page_listings = append(page_listings, page_listing)
	}

	template_err := templates.ExecuteTemplate(w,
		"index.html", page_listings)
	if template_err != nil {
		http.Error(w, template_err.Error(), http.StatusInternalServerError)
	}

	fmt.Println("pageIndexHandler end")
}

// as-is renders the add new page template
func addNewPageHandler(w http.ResponseWriter, r *http.Request) {
	template_err := templates.ExecuteTemplate(w, "add_new_page.html", nil)
	if template_err != nil {
		http.Error(w, template_err.Error(), http.StatusInternalServerError)
	}
}

// called by the add new page form
// makes a new blank page
func internalAddNewPageHandler(w http.ResponseWriter, r *http.Request) {
	title := r.FormValue("title")
	p := &Page{Title: title}
	err := p.save()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	filename := titleToFilename(title)
	http.Redirect(w, r, "/edit/"+filename, http.StatusFound)
}

// /////////////	Page View / Edit / Save Handlers //////////

var validPath = regexp.MustCompile("^/(edit|save|view)/([a-zA-Z0-9]+)$")

// meta handler for the page specific handlers
func makeHandler(fn func(http.ResponseWriter, *http.Request, string)) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		m := validPath.FindStringSubmatch(r.URL.Path)
		if m == nil {
			http.NotFound(w, r)
			return
		}
		fmt.Println(m[2])
		fn(w, r, m[2])
	}
}

// used via meta handler
func viewHandler(w http.ResponseWriter, r *http.Request, title string) {
	fmt.Println("viewHandler start")
	p, err := loadPage(title)
	if err != nil {
		http.Redirect(w, r, "/edit/"+title, http.StatusFound)
		return
	}
	template_err := templates.ExecuteTemplate(w, "view.html", p)
	if template_err != nil {
		http.Error(w, template_err.Error(), http.StatusInternalServerError)
	}
	fmt.Println("viewHandler end")
}

// used via meta handler
func editHandler(w http.ResponseWriter, r *http.Request, title string) {
	fmt.Println("editHandler start")
	p, err := loadPage(title)
	if err != nil {
		p = &Page{Title: title}
	}
	template_err := templates.ExecuteTemplate(w, "edit.html", p)
	if template_err != nil {
		http.Error(w, template_err.Error(), http.StatusInternalServerError)
	}
	fmt.Println("editHandler end")
}

// used via meta handler
func saveHandler(w http.ResponseWriter, r *http.Request, title string) {

	body := r.FormValue("body")
	p := &Page{Title: title, Body: []byte(body)}
	err := p.save()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/view/"+title, http.StatusFound)
}
