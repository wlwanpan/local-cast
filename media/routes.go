package media

import (
	"github.com/gorilla/mux"
	"github.com/wlwanpan/localcast/device"
)

func InitRoutes(r *mux.Router) *mux.Router {
	// Media endpoints
	r.HandleFunc("/media", GetMedia).Methods("GET")
	r.HandleFunc("/media/stop", device.UsesDevice(StopMedia)).Methods("POST")
	r.HandleFunc("/media/{id}/cast", device.UsesDevice(CastMedia)).Methods("POST")

	return r
}
