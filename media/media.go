package media

import (
	"io/ioutil"
	"path/filepath"
	"strings"

	"gopkg.in/mgo.v2/bson"
)

var (
	cachedMedia = map[string]*Media{}
)

type Media struct {
	ID   bson.ObjectId `bson:"_id"`
	Name string        `bson:"name"`
	Path string        `bson:"path"`
}

func (m *Media) GetID() string {
	return m.ID.Hex()
}

func NewMedia(name string, path string) *Media {
	return &Media{
		ID:   bson.NewObjectId(),
		Name: name,
		Path: path,
	}
}

func GetCachedMedia(mid string) (*Media, error) {
	m, ok := cachedMedia[mid]
	if ok {
		return m, nil
	}
	return &Media{}, ErrMediaNotFound
}

func GetAllCachedMedia() []*Media {
	m := []*Media{}
	for _, media := range cachedMedia {
		m = append(m, media)
	}
	return m
}

func LoadLocalFiles(p string) error {
	files, err := ioutil.ReadDir(p)
	if err != nil {
		return err
	}

	for _, file := range files {
		fileName := file.Name()
		if isHidden(fileName) {
			continue
		}
		newMedia := NewMedia(fileName, filepath.Join(p, fileName))
		cachedMedia[newMedia.GetID()] = newMedia
	}
	return nil
}

func CachedMediaCount() int {
	return len(cachedMedia)
}

func isHidden(filename string) bool {
	return strings.HasPrefix(filename, ".")
}
