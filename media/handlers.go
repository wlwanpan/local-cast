package media

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/context"
	"github.com/gorilla/mux"
	"github.com/wlwanpan/localcast/device"
)

type Payload struct {
	Data []*Media `json:"data"`
}

func GetMedia(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()

	var mt mediaType
	if query.Get("type") == "audio" {
		mt = AudioType
	} else {
		mt = VideoType
	}

	cachedMedias := Filter(mt, query.Get("search"))
	payload := &Payload{Data: cachedMedias}
	parsedPayload, err := json.Marshal(payload)
	if err != nil {
		log.Println(err)
	}

	w.WriteHeader(http.StatusOK)
	w.Write(parsedPayload)
}

func CastMedia(w http.ResponseWriter, r *http.Request) {
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

	device := context.Get(r, device.DeviceCtx).(device.Device)
	device.Start()

	// Fishy here but works better.
	log.Printf("Casting %s", media.Name)
	go func() {
		if err = device.PlayMedia(media.GetPath()); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		device.StopMedia()
	}()

	w.WriteHeader(http.StatusOK)
}

func StopMedia(w http.ResponseWriter, r *http.Request) {
	device := context.Get(r, device.DeviceCtx).(device.Device)
	if err := device.StopMedia(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
