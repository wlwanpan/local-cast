package chromecast

import (
	"errors"
	"log"

	"github.com/vishen/go-chromecast/application"
	"github.com/vishen/go-chromecast/dns"
)

const (
	// Chromecast is the local discovery name for chromecast device.
	Chromecast = "Home TV"

	// GoogleHome is the local discovery name for google home device.
	GoogleHome = "Home speaker"
)

var (
	// ErrNoChromecastAvailable is returned when attempting to use unavailable
	// chromecast device as an entry point
	ErrNoChromecastAvailable = errors.New("no chromecast connected to local network")
)

func ChromeApp() (*application.Application, error) {
	debug, disableCache := true, true
	app := application.NewApplication(debug, disableCache)
	entry, err := ChromecastDNSEntry()
	if err != nil {
		return &application.Application{}, err
	}

	if err := app.Start(entry); err != nil {
		return &application.Application{}, err
	}
	return app, nil
}

func ChromecastDNSEntry() (dns.CastEntry, error) {
	entries := dns.FindCastDNSEntries()
	targetDevice := GoogleHome
	if len(entries) == 0 {
		return dns.CastEntry{}, ErrNoChromecastAvailable
	}
	for _, entry := range entries {
		if entry.DeviceName == targetDevice {
			log.Printf("%s available.", targetDevice)
			return entry, nil
		}
	}
	return dns.CastEntry{}, ErrNoChromecastAvailable
}

func Play(p string) error {

	// path, content-type, transcode (.mp4)
	app, err := ChromeApp()
	if err != nil {
		return err
	}
	return app.Load(p, "", false)
}

func Stop() error {
	app, err := ChromeApp()
	if err != nil {
		return err
	}
	return app.Stop()
}
