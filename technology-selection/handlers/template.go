package handlers

import (
	"html/template"
	"net/http"
)

func Template_View(w http.ResponseWriter, data any, templateName string, files ...string) error {
	commonFiles := []string{
		"templates/common/header.gohtml",
		"templates/common/style.gohtml",
		"templates/common/styleFront.gohtml",
		"templates/common/checklist.gohtml",
	}
	files = append(files, commonFiles...)
	tmpl, err := template.New(templateName).ParseFiles(files...)
	if err != nil {
		return err
	}
	return tmpl.Execute(w, data)
}
