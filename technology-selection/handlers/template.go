package handlers

import (
	"html/template"
	"net/http"
)

func Template_View(w http.ResponseWriter, data any, templateName string, files ...string) error {
	tmpl, err := template.New(templateName).ParseFiles(files...)
	if err != nil {
		return err
	}
	return tmpl.Execute(w, data)
}
