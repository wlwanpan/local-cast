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
		Data: Load(),
	}
	data, err := json.Marshal(payload)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

// Err here -> once Stop cannot restart device ???
func StopDevice(w http.ResponseWriter, r *http.Request) {
	deviceUUID := r.Header.Get("device-uuid")
	device, err := GetByUUID(deviceUUID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := device.Stop(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
