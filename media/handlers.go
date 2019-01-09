package media

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/context"
	"github.com/gorilla/mux"
	"github.com/wlwanpan/localcast/device"
)

func GetHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()

	var mt mediaType
	if query.Get("type") == "audio" {
		mt = AudioType
	} else {
		mt = VideoType
	}

	cachedMedias := Filter(mt, query.Get("search"))
	payload := &struct {
		Data []*Media `json:"data"`
	}{
		cachedMedias,
	}
	parsedPayload, err := json.Marshal(payload)
	if err != nil {
		log.Println(err)
	}

	w.WriteHeader(http.StatusOK)
	w.Write(parsedPayload)
}

func CastHandler(w http.ResponseWriter, r *http.Request) {
	mid := mux.Vars(r)["id"]
	if mid == "" {
		log.Println("No id param provided.")
		http.Error(w, ErrMediaNotFound.Error(), http.StatusNotFound)
		return
	}
	media, err := Find(mid)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	device := context.Get(r, device.DeviceCtx).(*device.Device)

	// Fishy here but works better.
	log.Printf("Casting %s", media.Name)
	go func() {
		if err = device.PlayMedia(media.GetPath()); err != nil {
			log.Println(err)
		}
		device.StopMedia()
	}()

	w.WriteHeader(http.StatusOK)
}

func StopHandler(w http.ResponseWriter, r *http.Request) {
	device := context.Get(r, device.DeviceCtx).(*device.Device)
	if err := device.StopMedia(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func PauseHandler(w http.ResponseWriter, r *http.Request) {
	device := context.Get(r, device.DeviceCtx).(*device.Device)
	if err := device.PauseMedia(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func UnpauseHandler(w http.ResponseWriter, r *http.Request) {
	device := context.Get(r, device.DeviceCtx).(*device.Device)
	if err := device.UnpauseMedia(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
