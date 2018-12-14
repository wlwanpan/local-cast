package device

import "github.com/gorilla/mux"

func InitRoutes(r *mux.Router) *mux.Router {
	r.HandleFunc("/device/stop", StopDevice).Methods("POST")
	return r
}
