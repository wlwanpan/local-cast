package main

import (
	"log"
	"net/http"
	"time"

	"github.com/wlwanpan/localcast/media"
)

func main() {
	mediaPath := "/Users/warrenwan/Music"
	log.Printf("Loading localfiles from: %s \n", mediaPath)
	if err := media.LoadLocalFiles(mediaPath); err != nil {
		log.Fatal(err)
	}
	log.Printf("Sucessfully cached: %d", media.CachedMediaCount())

	port := ":4040"
	log.Printf("Starting server on %s", port)
	s := &http.Server{
		Handler:      media.InitMediaRoutes(),
		Addr:         port,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	log.Fatal(s.ListenAndServe())
}
