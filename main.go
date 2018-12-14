package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/akamensky/argparse"
	"github.com/gorilla/mux"
	"github.com/wlwanpan/localcast/chromecast"
	"github.com/wlwanpan/localcast/common"
	"github.com/wlwanpan/localcast/device"
	"github.com/wlwanpan/localcast/media"
)

func initRoutes() http.Handler {
	r := mux.NewRouter()
	r = media.InitRoutes(r)
	r = device.InitRoutes(r)
	return common.SetCORS(r)
}

func main() {
	parser := argparse.NewParser("localcast", "Api server for casting local media files.")
	log.SetOutput(os.Stdout)

	port := parser.String("p", "port", &argparse.Options{
		Required: false,
		Default:  "4040",
		Help:     "Port to run the server.",
	})
	mediaSrc := parser.String("s", "source", &argparse.Options{
		Required: true,
		Help:     "Path of folder to serve.",
	})

	if err := parser.Parse(os.Args); err != nil {
		log.Fatal(err)
	}

	log.Printf("Loading localfiles from: %s \n", *mediaSrc)
	if err := media.LoadLocalFiles(*mediaSrc); err != nil {
		log.Fatal(err)
	}
	log.Printf("Sucessfully cached: %d", media.CachedMediaCount())

	log.Printf("Initializing Google Home")
	if err := chromecast.InitGoogleHomeApp(); err != nil {
		log.Fatal(err)
	}

	s := &http.Server{
		Handler:      initRoutes(),
		Addr:         ":" + *port,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	log.Printf("Starting server on %s\n", *port)
	log.Fatal(s.ListenAndServe())
}
