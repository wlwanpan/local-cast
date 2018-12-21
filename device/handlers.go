package device

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/context"
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
	device := context.Get(r, DeviceCtx).(Device)
	if err := device.Stop(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
