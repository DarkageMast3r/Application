package handlers

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"

	"main/models"
)

type Service struct {
	Hosts      []string
	LastServed int
}

func Service_Register(w http.ResponseWriter, r *http.Request) {
	idx := strings.LastIndex(r.RemoteAddr, ":")
	serviceUri := r.RemoteAddr[:idx] + ":" + r.PathValue("port")
	models.Service_Create(r.PathValue("service"), serviceUri)
}

func Service_Get_Names(w http.ResponseWriter, r *http.Request) {
	json, err := json.Marshal(models.Service_Get_Names())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(json)
}

func Service_Get(w http.ResponseWriter, r *http.Request) {
	service, exists := models.Service_Get_By_Name(r.PathValue("service"))
	if !exists {
		http.NotFound(w, r)
		return
	}
	service.LastServed = (service.LastServed + 1) % len(service.Hosts)
	service.Save()
	io.WriteString(w, service.Hosts[service.LastServed])
}
