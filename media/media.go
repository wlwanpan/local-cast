package media

import (
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/h2non/filetype"
	"github.com/lithammer/fuzzysearch/fuzzy"
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
	return filepath.Join(m.path, m.Name) + m.extension
}

func New(name, extension, path string, mediaType mediaType) *Media {
	return &Media{
		ID:        bson.NewObjectId(),
		Name:      name,
		extension: extension,
		path:      path,
		mediaType: mediaType,
	}
}

func Find(mid string) (*Media, error) {
	m, ok := cachedMedia[mid]
	if ok {
		return m, nil
	}
	return &Media{}, ErrMediaNotFound
}

func Filter(mt mediaType, search string) []*Media {
	m := []*Media{}
	for _, media := range cachedMedia {
		if media.mediaType != mt {
			continue
		}
		if search != "" {
			s := strings.ToLower(search)
			t := strings.ToLower(media.Name)
			if !fuzzy.Match(s, t) {
				continue
			}
		}
		m = append(m, media)
	}
	return m
}

// LoadLocalFiles recursively reads paths and cache the media files.
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
			return LoadLocalFiles(filePath)
		}
		fileType := readFileType(filePath)
		if fileType != UnknownType {
			name, extension := fileExtension(fileName)
			newMedia := New(name, extension, p, fileType)
			cachedMedia[newMedia.GetID()] = newMedia
		}
	}
	return nil
}

// CachedMediaCount get the amount of media cached in mem.
func CachedMediaCount() int {
	return len(cachedMedia)
}

func isHidden(filename string) bool {
	return strings.HasPrefix(filename, ".")
}

// fileExtension extracts the file extension from the filename.
func fileExtension(filename string) (string, string) {
	i := strings.LastIndex(filename, ".")
	return filename[:i], filename[i:]
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
