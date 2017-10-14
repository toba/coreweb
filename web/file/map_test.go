package file_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"toba.tech/app/lib/web/file"
	"toba.tech/app/lib/web/header/content"
	"toba.tech/app/lib/web/mime"
)

func TestMapRead(t *testing.T) {
	m, err := file.InFolder(folder, true)
	assert.NoError(t, err)
	err = m.Read(false)
	assert.NoError(t, err)

	file0a := "file0a.txt"
	file2b := "folder2" + slash + "file2b.txt"

	assert.Contains(t, m.Files, file0a)
	assert.Contains(t, m.Files, file2b)

	text := string(m.Files[file2b].Content)
	assert.Equal(t, "file2b", text)

	head := m.Files[file0a].Header
	assert.Equal(t, "6", head[content.Length])
	assert.Equal(t, mime.Text, head[content.Type])
	//assert.Equal(t, "Mon, 06 Mar 2017 21:04:14 MST", head[header.LastModified])
}
