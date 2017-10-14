package file

import (
	"archive/zip"
	"errors"
	"fmt"
	"io/ioutil"
	"strings"
)

var (
	zipData      string
	ErrNoZipData = errors.New("No zip data found")
)

// HasZipData indicates if a zip string exists without validating it.
func HasZipData() bool {
	return zipData != ""
}

// RegisterZip assigns zip data containing embedded files. This is usually
// called by generated code.
func RegisterZip(data string) {
	zipData = data
}

// InZipFile returns all files inside the zip file data. Based on
// https://github.com/rakyll/statik
func InZipFile() (*Map, error) {
	if zipData == "" {
		return nil, ErrNoZipData
	}
	zipReader, err := zip.NewReader(strings.NewReader(zipData), int64(len(zipData)))
	if err != nil {
		return nil, err
	}
	files := &Map{Files: make(map[string]*Info)}

	for _, zipFile := range zipReader.File {
		unzipped, err := unzip(zipFile)
		if err != nil {
			return nil, fmt.Errorf("error unzipping file %q: %s", zipFile.Name, err)
		}
		files.addZip(unzipped, zipFile)
	}
	return files, nil
}

func unzip(zf *zip.File) ([]byte, error) {
	rc, err := zf.Open()
	if err != nil {
		return nil, err
	}
	defer rc.Close()
	return ioutil.ReadAll(rc)
}
