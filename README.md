# Localcast

An API server built around [go-chromecast](https://github.com/vishen/go-chromecast) to cast your local media files to
[Google Home](https://store.google.com/ca/product/google_home) or [Chromecast](https://store.google.com/product/chromecast).

Note: The server has to be under the same local network.

# Start Server
```
go run main.go -s PATH_TO_DIR
```

# Endpoints

- Get list of media in local dir:
```bash
curl localhost:4040/media?type=audio
```

Response:
```json
[
  {
    "ID":"5c0dbd824f3d37099f988541",
    "Name":"Running.mp3",
  },
  {
    "ID":"5c0dbd824f3d37099f988542",
    "Name":"Walking.mp3",
  }
]
```

- Cast media by id:
```bash
curl -X POST localhost:4040/media/5c0dbd824f3d37099f988541/cast
```

- Stop media:
```bash
curl -X POST localhost:4040/media/stop
```