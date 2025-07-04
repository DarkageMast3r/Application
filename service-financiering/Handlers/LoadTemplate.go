package handlers

import (
	"fmt"
	tmpl "html/template"
	"net/http"
)

func LoadTemplate(wr http.ResponseWriter, templateName string, data any) error {
	template, err := tmpl.New("Home").ParseFiles(templateName)
	if err != nil {
		fmt.Println("error rendering template")
		return err
	}
	return template.Execute(wr, data)
}
