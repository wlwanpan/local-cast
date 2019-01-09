package device

import (
	"errors"
	"log"

	gcapp "github.com/vishen/go-chromecast/application"
	"github.com/vishen/go-chromecast/dns"
)

var (
	cachedDevices = []*Device{}
)

type Device struct {
	Name      string `json:"name"`
	UUID      string `json:"uuid"`
	app       *gcapp.Application
	castEntry *dns.CastEntry
}

func NewDevice(entry *dns.CastEntry) *Device {
	return &Device{
		Name:      entry.GetName(),
		UUID:      entry.GetUUID(),
		app:       gcapp.NewApplication(false, false), // debug, disableCache
		castEntry: entry,
	}
}

func (d *Device) Start() error {
	return d.app.Start(d.castEntry)
}

func (d *Device) Stop() error {
	return d.app.Stop()
}

func (d *Device) Close() {
	d.app.Close()
}

func (d *Device) StopMedia() error {
	return d.app.StopMedia()
}

func (d *Device) PauseMedia() error {
	_, castMedia, _ := d.app.Status()
	if castMedia == nil {
		log.Println("No media currently playing")
		return errors.New("no media playing")
	}
	return d.app.Pause()
}

func (d *Device) UnpauseMedia() error {
	return d.app.Unpause()
}

func (d *Device) PlayMedia(p string) error {
	_, castMedia, _ := d.app.Status()
	if castMedia != nil {
		log.Println("Stopping current media ...")
		if err := d.StopMedia(); err != nil {
			return err
		}
		return d.app.Load(p, "", false)
	}
	d.Start()
	if err := d.app.Load(p, "", false); err != nil {
		return err
	}
	return nil
}

func GetByUUID(uuid string) (*Device, error) {
	for _, d := range cachedDevices {
		if d.UUID == uuid {
			return d, nil
		}
	}
	return &Device{}, ErrDeviceNotFound
}

func LoadAndCache() []*Device {
	devices := Load()
	Cache(devices)
	return devices
}

func Cache(devices []*Device) {
	if len(cachedDevices) != 0 {
		for _, cachedDevice := range cachedDevices {
			cachedDevice.Stop()
		}
	}
	cachedDevices = devices
}

func Load() []*Device {
	entries := dns.FindCastDNSEntries()
	devices := []*Device{}
	for _, entry := range entries {
		devices = append(devices, NewDevice(&entry))
	}
	return devices
}
