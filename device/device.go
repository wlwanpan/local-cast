package device

import (
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
	if err := d.app.Stop(); err != nil {
		return err
	}
	// d.app.Close()
	return nil
}

func (d *Device) StopMedia() error {
	return d.app.StopMedia()
}

func (d *Device) PlayMedia(p string) error {
	_, castMedia, _ := d.app.Status()
	if castMedia != nil {
		if err := d.StopMedia(); err != nil {
			return err
		}
	}
	return d.app.Load(p, "", false)
}

func GetByUUID(uuid string) (*Device, error) {
	for _, d := range cachedDevices {
		if d.UUID == uuid {
			return d, nil
		}
	}
	return &Device{}, ErrDeviceNotFound
}

func Cache(devices []*Device) {
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
