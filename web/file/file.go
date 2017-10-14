// Package file manages files for HTTP service.
package file

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"strconv"

	"time"

	"toba.tech/app/lib/web/encoding"
	"toba.tech/app/lib/web/header"
	"toba.tech/app/lib/web/header/content"
	"toba.tech/app/lib/web/mime"
)

const slash = string(os.PathSeparator)

var (
	wd         = ""
	normalized = false
	resolver   = func() (string, error) {
		ex, err := os.Executable()
		if err != nil {
			return "", err
		}
		return filepath.Dir(ex), nil
	}
)

// Resolve allows test methods to specify a custom working directory resolver.
func Resolve(r func() (string, error)) {
	resolver = r
}

// normalize appends an OS slash as needed to the root working directory.
func normalize() {
	if !normalized {
		if wd == "" {
			path, err := resolver()
			if err != nil {
				log.Fatal(err)
			}
			wd = path
		}
		if wd != "" && !strings.HasSuffix(wd, slash) {
			wd += slash
		}
		log.Printf("Set working directory to %s", wd)
		normalized = true
	}
}

// Open returns a file handle for a named file in the working directory.
func Open(fileName string) (*os.File, error) {
	normalize()
	return os.Open(wd + fileName)
}

// Read returns bytes for a named file in the working directory.
func Read(fileName string) ([]byte, error) {
	normalize()
	return ioutil.ReadFile(wd + fileName)
}

// InFolder retrieves all file paths in a directory with the option to also
// retrieve file paths from sub-directories.
func InFolder(path string, recursive bool) (*Map, error) {
	normalize()
	folder := wd + path
	files := &Map{Files: make(map[string]*Info)}

	if recursive {
		return subFolderFiles(files, folder, folder)
	} else {
		info, err := ioutil.ReadDir(folder)
		if err != nil {
			return nil, err
		}
		for _, f := range info {
			if !f.IsDir() {
				files.add(folder+slash+f.Name(), folder, f)
			}
		}
		return files, nil
	}
}

// makeHeader writes the HTTP header values for the file.
func makeHeader(f os.FileInfo) map[string]string {
	h := map[string]string{
		content.Type:        mime.Infer(f.Name()),
		content.Length:      strconv.FormatInt(f.Size(), 10),
		header.LastModified: f.ModTime().Format(time.RFC1123),
	}
	if strings.HasSuffix(f.Name(), ".gz") {
		h[content.Encoding] = encoding.GZip
	}
	return h
}

// makeRelative removes the absolute working directory from a path.
func makeRelative(filePath, rootPath string) string {
	return strings.Replace(filePath, rootPath+slash, "", -1)
}

// subFolderFiles visits all subfolders recursively to generate a resursive
// file map.
func subFolderFiles(m *Map, rootFolder, subFolder string) (*Map, error) {
	info, err := ioutil.ReadDir(subFolder)

	if err != nil {
		return nil, err
	}

	for _, f := range info {
		subPath := subFolder + slash + f.Name()
		if f.IsDir() {
			subFiles, err := subFolderFiles(m, rootFolder, subPath)
			if err != nil {
				return nil, err
			}
			m.addAll(subFiles)
		} else {
			m.add(subPath, rootFolder, f)
		}
	}

	return m, nil
}

// Monitor files and update map with changed content. This is only active
// during debug so use simple polling rather than hooking OS notify:
// https://github.com/fsnotify/fsnotify
func Monitor(m *Map) {
	ticker := time.NewTicker(3 * time.Second)
	go func() {
		for _ = range ticker.C {
			err := UpdateChangedFiles(m)
			if err != nil {
				println(err.Error())
				ticker.Stop()
			}
		}
	}()
}

// UpdateChangedFiles updates file content in map if the file system modified
// time is newer.
func UpdateChangedFiles(m *Map) error {
	for _, info := range m.Files {
		f, err := os.Stat(info.Path)
		if err != nil {
			println("error getting info for " + info.Path)
			return err
		}
		if f.ModTime().After(info.Modified) {
			content, err := ioutil.ReadFile(info.Path)
			if err != nil {
				println("failed reloading " + info.Path)
				return err
			}
			println("detected change in " + info.Path)

			m.Lock()
			info.Content = content
			info.Modified = f.ModTime()
			info.Compressed = nil
			info.Compress()
			m.Unlock()
		}
	}
	return nil
}
