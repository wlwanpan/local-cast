package device

import (
	"log"
	"net/http"

	"github.com/wlwanpan/localcast/chromecast"
)

func StopDevice(w http.ResponseWriter, r *http.Request) {
	if err := chromecast.StopGoogleHomeApp(); err != nil {
		log.Println(err)
	}
	w.WriteHeader(http.StatusOK)
}
