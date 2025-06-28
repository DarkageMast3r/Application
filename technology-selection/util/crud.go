package util

import (
	"fmt"
	"net/http"
	"reflect"
)

func Fill_Fields_From_Form[k any](object *k, r *http.Request) error {
	objValue := reflect.ValueOf(object).Elem()
	objType := reflect.TypeOf(*object)

	fieldCount := objType.NumField()
	for i := range fieldCount {
		field := objType.Field(i)
		jsonName, exists := field.Tag.Lookup("json")
		_, shouldExclude := field.Tag.Lookup("excludeFromCreate")
		if !exists || shouldExclude {
			continue
		}

		if !r.Form.Has(jsonName) {
			return http.ErrBodyNotAllowed
		}

		val := objValue.Field(i)
		formVal := r.Form.Get(jsonName)
		if val.Kind() == reflect.String {
			if val.CanSet() {
				val.SetString(formVal)
			} else {
				fmt.Println("Cannot set field ", objType.Name(), "/", field.Name)
			}
		}
	}
	return nil
}

func Crud_View_Create(w http.ResponseWriter, objType reflect.Type, action string) {
	w.WriteHeader(http.StatusOK)
	w.Header().Add("content-type", "text/html")
	w.Write([]byte("<html><head><title>Create</title></head><body>"))
	w.Write([]byte(fmt.Sprintf("<form action=\"%s\" method=\"post\">", action)))
	w.Write([]byte("<div style=\"display: grid; grid-template-columns: 1fr 1fr; gap: 0.5rem 1rem\">"))

	fieldCount := objType.NumField()
	for i := range fieldCount {
		field := objType.Field(i)
		jsonName, exists := field.Tag.Lookup("json")
		_, shouldExclude := field.Tag.Lookup("excludeFromCreate")
		if exists && !shouldExclude {
			text := fmt.Sprintf(
				"<label style=\"text-align: end\" for=\"%s\">%s</label><input id=\"%s\" name=\"%s\" type=\"text\"/>",
				jsonName,
				field.Name,
				jsonName,
				jsonName,
			)
			w.Write([]byte(text))
		}
	}
	w.Write([]byte("<input type=\"submit\"/>"))
	w.Write([]byte("</div></form>"))
	w.Write([]byte("</body>"))
}
