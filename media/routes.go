package media

import (
	"github.com/gorilla/mux"
)

func InitRoutes(r *mux.Router) *mux.Router {
	// Media endpoints
	r.HandleFunc("/media", GetMedia).Methods("GET")
	r.HandleFunc("/media/stop", StopMedia).Methods("POST")
	r.HandleFunc("/media/{id}/cast", CastMedia).Methods("POST")

	return r
}
