package device

import (
	"net/http"

	"github.com/gorilla/context"
)

const (
	DeviceCtx = "DEVICE_CTX"
)

func UsesDevice(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		uuid := r.Header.Get("device-uuid")
		device, err := GetByUUID(uuid)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError) // Add proper status error
			return
		}
		context.Set(r, DeviceCtx, device)
		fn(w, r)
	}
}
