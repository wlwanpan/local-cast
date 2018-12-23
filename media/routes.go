package media

import (
	"github.com/gorilla/mux"
	"github.com/wlwanpan/localcast/device"
)

func InitRoutes(r *mux.Router) *mux.Router {
	// Media endpoints
	r.HandleFunc("/media", GetHandler).Methods("GET")
	r.HandleFunc("/media/{id}/cast", device.UsesDevice(CastHandler)).Methods("POST")
	r.HandleFunc("/media/stop", device.UsesDevice(StopHandler)).Methods("POST")

	return r
}
