package media

import (
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/h2non/filetype"
	"gopkg.in/mgo.v2/bson"
)

const (
	AudioType   = iota
	VideoType   = iota
	UnknownType = iota
)

type mediaType int

var (
	cachedMedia = map[string]*Media{}
)

type Media struct {
	ID        bson.ObjectId `json:"_id"`
	Name      string        `json:"name"`
	extension string
	path      string
	mediaType mediaType
}

func (m *Media) GetID() string {
	return m.ID.Hex()
}

func (m *Media) GetPath() string {
	return m.path
}

func NewMedia(name string, path string, mediaType mediaType) *Media {
	return &Media{
		ID:        bson.NewObjectId(),
		Name:      name,
		path:      path,
		mediaType: mediaType,
	}
}

func GetCachedMedia(mid string) (*Media, error) {
	m, ok := cachedMedia[mid]
	if ok {
		return m, nil
	}
	return &Media{}, ErrMediaNotFound
}

func GetCachedAudio() []*Media {
	return filterCachedMedia(AudioType)
}

func GetCachedVideo() []*Media {
	return filterCachedMedia(VideoType)
}

func LoadLocalFiles(p string) error {
	files, err := ioutil.ReadDir(p)
	if err != nil {
		return err
	}

	for _, file := range files {
		fileName := file.Name()
		filePath := filepath.Join(p, fileName)
		if isHidden(fileName) {
			continue
		}
		if file.IsDir() {
			LoadLocalFiles(filePath)
		}
		fileType := readFileType(filePath)
		if fileType != UnknownType {
			newMedia := NewMedia(fileName, filePath, fileType)
			cachedMedia[newMedia.GetID()] = newMedia
		}
	}
	return nil
}

func CachedMediaCount() int {
	return len(cachedMedia)
}

func isHidden(filename string) bool {
	return strings.HasPrefix(filename, ".")
}

func filterCachedMedia(mt mediaType) []*Media {
	m := []*Media{}
	for _, media := range cachedMedia {
		if media.mediaType == mt {
			m = append(m, media)
		}
	}
	return m
}

func readFileType(f string) mediaType {
	buf, _ := ioutil.ReadFile(f)
	if filetype.IsAudio(buf) {
		return AudioType
	}
	if filetype.IsVideo(buf) {
		return VideoType
	}
	return UnknownType
}
