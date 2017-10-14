package file_test

import (
	"os"
	"testing"

	"toba.tech/app/lib/web/file"

	"github.com/stretchr/testify/assert"
)

const (
	slash  = string(os.PathSeparator)
	folder = "lib" + slash + "web" + slash + "file" + slash + "test"
)

// TestMain applies a custom working directory resolver since tests may run
// in a temporary directory.
func TestMain(m *testing.M) {
	code := 1
	defer func() {
		if code == 0 {
			code = m.Run()
		}
		os.Exit(code)
	}()

	file.Resolve(func() (string, error) {
		return os.Getenv("TOBA_PATH"), nil
	})

	code = 0
}

func TestFileOpen(t *testing.T) {
	f, err := file.Open("toba.config")
	assert.NoError(t, err)
	assert.NotNil(t, f)

	_, err = file.Open("nosuchfile.txt")
	assert.Error(t, err)
}

func TestFileRead(t *testing.T) {
	data, err := file.Read("toba.config")
	assert.NoError(t, err)
	assert.NotNil(t, data)

	text := string(data)
	assert.Contains(t, text, "sslCert")
}

func TestFileInFolder(t *testing.T) {
	m, err := file.InFolder(folder, false)
	assert.NoError(t, err)
	assert.Len(t, m.Files, 2)

	m, err = file.InFolder(folder, true)
	assert.NoError(t, err)
	assert.Len(t, m.Files, 6)
}

func writeFile(t *testing.T, f *os.File, content string) {
	err := f.Truncate(0)
	assert.NoError(t, err)

	_, err = f.WriteString(content)
	assert.NoError(t, err)

	err = f.Close()
	assert.NoError(t, err)
}

func TestFileUpdate(t *testing.T) {
	wd, err := os.Getwd()
	assert.NoError(t, err)

	fileName := "update.txt"
	testPath := wd + slash + "test" + slash + fileName
	before := "before content"
	after := "after content"

	temp, err := os.Create(testPath)
	assert.NoError(t, err)

	defer os.Remove(testPath)

	writeFile(t, temp, before)

	m, err := file.InFolder(folder, false)
	assert.NoError(t, err)
	assert.Len(t, m.Files, 3)
	assert.Contains(t, m.Files, fileName)

	err = m.Read(false)
	assert.NoError(t, err)

	info := m.Files[fileName]
	assert.Equal(t, before, string(info.Content))

	temp, err = os.OpenFile(testPath, os.O_WRONLY, os.ModeAppend)
	assert.NoError(t, err)

	writeFile(t, temp, after)

	err = file.UpdateChangedFiles(m)
	assert.NoError(t, err)

	assert.Equal(t, after, string(info.Content))
	assert.NotNil(t, info.Compressed)

	// err = os.Remove(testPath)
	// assert.NoError(t, err)
}
