package media

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/dhowden/tag"
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
	Metadata  *Metadata     `json:"metadata"`
	extension string
	path      string
	mediaType mediaType
}

type Metadata struct {
	Title     string `json:"title"`
	Artist    string `json:"artist"`
	Album     string `json:"album"`
	Thumbnail string `json:"thumbnail"`
}

func (m *Media) GetID() string {
	return m.ID.Hex()
}

func (m *Media) GetPath() string {
	return filepath.Join(m.path, m.Name) + m.extension
}

func New(name, extension, path string, mediaType mediaType, md *Metadata) *Media {
	return &Media{
		ID:        bson.NewObjectId(),
		Name:      name,
		Metadata:  md,
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
		fileType, md, err := getMediaType(filePath)
		if err != nil {
			log.Printf("Could not get metadata of %s", filePath)
		}
		if fileType != UnknownType {
			name, extension := fileExtension(fileName)
			newMedia := New(name, extension, p, fileType, md)
			cachedMedia[newMedia.GetID()] = newMedia
		}
	}
	return nil
}

func getMediaType(p string) (mediaType, *Metadata, error) {
	f, _ := os.Open(p)
	head := make([]byte, 261)
	f.Read(head)
	if filetype.IsAudio(head) {
		md, err := getAudioMetadata(f)
		return AudioType, md, err
	}
	if filetype.IsVideo(head) {
		return VideoType, &Metadata{}, nil
	}
	return UnknownType, &Metadata{}, nil
}

func getAudioMetadata(file *os.File) (*Metadata, error) {
	m, err := tag.ReadFrom(file)
	if err != nil {
		return &Metadata{}, err
	}
	var thumbnail string
	if pic := m.Picture(); pic != nil {
		log.Println(pic)
		thumbnail = pic.String()
	}
	return &Metadata{
		Title:     m.Title(),
		Artist:    m.Artist(),
		Album:     m.Album(),
		Thumbnail: thumbnail,
	}, nil
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
