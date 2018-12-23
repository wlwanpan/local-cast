package device

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/context"
)

type Payload struct {
	Data []*Device `json:"data"`
}

func GetHandler(w http.ResponseWriter, r *http.Request) {
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
func StopHandler(w http.ResponseWriter, r *http.Request) {
	device := context.Get(r, DeviceCtx).(*Device)
	if err := device.Stop(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func RefreshHandler(w http.ResponseWriter, r *http.Request) {
	devices := LoadAndCache()
	log.Printf("Cached %d devices", len(devices))
	w.WriteHeader(http.StatusOK)
}
