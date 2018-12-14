package chromecast

import (
	"errors"
	"log"

	"github.com/vishen/go-chromecast/application"
	"github.com/vishen/go-chromecast/dns"
)

const (
	// ChromecastName is the local discovery name for chromecast device.
	ChromecastName = "Home TV"

	// GoogleHomeName is the local discovery name for google home device.
	GoogleHomeName = "Home speaker"
)

var (
	// ErrNoChromecastAvailable is returned when attempting to use unavailable
	// chromecast device as an entry point
	ErrNoChromecastAvailable = errors.New("no chromecast connected to local network")
)

var (
	ccApp *application.Application
	ghApp *application.Application
)

func InitChromecastApp() error {
	var err error
	ccApp, err = initApp(ChromecastName)
	return err
}

func RefreshChromecastApp() error {
	if err := ccApp.Stop(); err != nil {
		return err
	}
	return InitChromecastApp()
}

func InitGoogleHomeApp() error {
	var err error
	ghApp, err = initApp(GoogleHomeName)
	return err
}

func RefreshGoogleHomeApp() error {
	if err := ghApp.Stop(); err != nil {
		return err
	}
	return InitGoogleHomeApp()
}

func StopGoogleHomeApp() error {
	return ghApp.Stop()
}

func initApp(targetDevice string) (*application.Application, error) {
	debug, disableCache := true, true
	app := application.NewApplication(debug, disableCache)
	entry, err := chromecastDNSEntry(targetDevice)
	if err != nil {
		return &application.Application{}, err
	}

	if err = app.Start(entry); err != nil {
		return &application.Application{}, err
	}

	return app, nil
}

func chromecastDNSEntry(targetDevice string) (dns.CastEntry, error) {
	entries := dns.FindCastDNSEntries()
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

func PlayMedia(p string) error {
	return ghApp.Load(p, "", false)
}

func StopMedia() error {
	return ghApp.StopMedia()
}
