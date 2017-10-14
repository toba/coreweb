package file

import (
	"archive/zip"
	"io/ioutil"
	"os"
	"sync"
)

// Map file path to information about the file.
type Map struct {
	sync.RWMutex
	Files map[string]*Info
}

// Read updates all Content bytes in the Map and optionally GZips them.
func (m *Map) Read(gzip bool) error {
	normalize()
	for _, info := range m.Files {
		if info.Content == nil {
			content, err := ioutil.ReadFile(info.Path)
			if err != nil {
				return err
			}
			info.Content = content
		}

		if gzip {
			err := info.Compress()
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// add os.FileInfo to the Map.
func (m *Map) add(filePath, rootFolder string, info os.FileInfo) {
	m.Files[makeRelative(filePath, rootFolder)] = &Info{
		Header:   makeHeader(info),
		Path:     filePath,
		Modified: info.ModTime(),
	}
}

// addAll adds all members of one file Map to this one.
func (m *Map) addAll(other *Map) {
	for k, v := range other.Files {
		m.Files[k] = v
	}
}

// addFromZip adds zip file information to the Map.
func (m *Map) addZip(content []byte, f *zip.File) {
	m.Files[f.Name] = &Info{
		Header:   makeHeader(f.FileInfo()),
		Content:  content,
		Path:     f.Name,
		Modified: f.ModTime(),
	}
}
