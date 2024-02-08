package main

import (
	"fmt"
	"os"
	"strings"
)

func titleToFilename(title string) string {
	return strings.ReplaceAll(title, " ", "_")
}

func filenameToTitle(filename string) string {
	return strings.ReplaceAll(filename, "_", " ")
}

// Saves a page struct to a file
// The file is named after the page's title, and the content is the just a dump of the page's body
func (p *Page) save() error {
	// build the filename
	trimmed := strings.TrimSpace(p.Title)
	filename := strings.ReplaceAll(trimmed, " ", "_")
	if _, err := os.Stat(saved_pages_path); os.IsNotExist(err) {
		os.Mkdir(saved_pages_path, 0755)
	}
	fullpath := saved_pages_path + "/" + filename
	return os.WriteFile(fullpath, p.Body, 0600)
}

// Loads a page from a file
// The file is named after the page's title, and the body is the file content as a byte array
func loadPage(title string) (*Page, error) {
	fmt.Println("attempt to load page with title: ", title)
	// invert the save process
	filename := strings.ReplaceAll(title, "_", " ")
	fullpath := saved_pages_path + "/" + filename
	// print the fullpath
	fmt.Println("fullpath: ", fullpath)
	body, err := os.ReadFile(fullpath)
	if err != nil {
		return nil, err
	}
	return &Page{Title: title, Body: body}, nil
}
