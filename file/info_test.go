package file_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/toba/coreweb/encoding"
	"github.com/toba/coreweb/file"
	"github.com/toba/coreweb/header/content"
	"github.com/toba/coreweb/mime"
)

func TestInfoReplace(t *testing.T) {
	m, err := file.InFolder(folder, true)
	assert.NoError(t, err)
	err = m.Read(false)
	assert.NoError(t, err)

	info := m.Files["folder1"+slash+"file1a.txt"]
	text := string(info.Content)
	assert.Equal(t, text, "file1a")

	updated := info.Replace("file", "replace")
	newText := string(updated.Content)
	assert.Equal(t, newText, "replace1a")
}

func TextInfoCompressible(t *testing.T) {
	info := &file.Info{
		Header: map[string]string{
			content.Type: mime.HTML,
		},
	}
	assert.False(t, info.Compressible())

	info.Content = []byte{1, 2, 3, 4, 5, 6, 7, 8}
	assert.True(t, info.Compressible())

	info.Header[content.Type] = mime.PNG
	assert.False(t, info.Compressible())

	info.Header[content.Type] = mime.Text
	info.Header[content.Encoding] = encoding.GZip
	assert.False(t, info.Compressible())
}

// TestInfoCompress ensures compressible file is GZipped and has GZip header
// but that parent info retains standard header fields.
func TestInfoCompress(t *testing.T) {
	m, err := file.InFolder(folder, true)
	assert.NoError(t, err)
	err = m.Read(true)
	assert.NoError(t, err)

	info := m.Files["folder1"+slash+"file1a.txt"]

	assert.NotNil(t, info.Compressed)
	assert.Equal(t, encoding.GZip, info.Compressed.Header[content.Encoding])
	assert.NotEqual(t, encoding.GZip, info.Header[content.Encoding])
}
