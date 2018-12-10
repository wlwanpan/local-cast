package media

import (
	"net/http"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func InitMediaRoutes() http.Handler {
	r := mux.NewRouter()

	// Media endpoints
	r.HandleFunc("/media", GetMedia).Methods("GET")
	r.HandleFunc("/stop", StopMedia).Methods("POST")
	r.HandleFunc("/media/{id}/cast", CastMedia).Methods("POST")

	headers := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type"})
	origins := handlers.AllowedOrigins([]string{"*"})
	methods := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "OPTIONS"})
	return handlers.CORS(headers, origins, methods)(r)
}
