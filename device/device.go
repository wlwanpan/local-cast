package device

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

type Device struct {
	Name      string `json:"name"`
	UUID      string `json:"uuid"`
	castEntry *dns.CastEntry
}

func NewDevice(entry *dns.CastEntry) *Device {
	return &Device{
		Name:      entry.GetName(),
		UUID:      entry.GetUUID(),
		castEntry: entry,
	}
}

func LoadDevices() []*Device {
	entries := dns.FindCastDNSEntries()
	devices := []*Device{}
	for _, entry := range entries {
		devices = append(devices, NewDevice(&entry))
	}
	return devices
}

func InitGoogleHomeApp() error {
	var err error
	ghApp, err = initApp(GoogleHomeName)
	return err
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
	_, castMedia, _ := ghApp.Status()
	if castMedia != nil {
		if err := StopMedia(); err != nil {
			return err
		}
	}
	return ghApp.Load(p, "", false)
}

func StopMedia() error {
	return ghApp.StopMedia()
}
