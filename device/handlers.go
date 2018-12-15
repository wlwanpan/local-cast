package device

import (
	"encoding/json"
	"net/http"
)

type Payload struct {
	Data []*Device `json:"data"`
}

func GetDevice(w http.ResponseWriter, r *http.Request) {
	payload := &Payload{
		Data: LoadDevices(),
	}
	data, err := json.Marshal(payload)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

func StopDevice(w http.ResponseWriter, r *http.Request) {
	if err := ghApp.Stop(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
