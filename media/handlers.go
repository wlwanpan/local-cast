package media

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/wlwanpan/localcast/chromecast"
)

type Payload struct {
	Data []*Media `json:"data"`
}

func GetMedia(w http.ResponseWriter, r *http.Request) {
	cachedMedias := GetAllCachedMedia()
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
	}
	media, err := GetCachedMedia(mid)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	log.Printf("Casting %s", media.Name)
	if err = chromecast.Play(media.Path); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusOK)
}

func StopMedia(w http.ResponseWriter, r *http.Request) {
	if err := chromecast.Stop(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.WriteHeader(http.StatusOK)
}
