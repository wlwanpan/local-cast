package device

import "github.com/gorilla/mux"

func InitRoutes(r *mux.Router) *mux.Router {
	// Device endpoints
	r.HandleFunc("/devices", GetHandler).Methods("GET")
	r.HandleFunc("/device/stop", UsesDevice(StopHandler)).Methods("POST")
	r.HandleFunc("/device/refresh", RefreshHandler).Methods("POST")
	return r
}
