package device

import "github.com/gorilla/mux"

func InitRoutes(r *mux.Router) *mux.Router {
	r.HandleFunc("/devices", GetDevice).Methods("GET")
	r.HandleFunc("/device/stop", StopDevice).Methods("POST")
	return r
}
