package mime_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/toba/coreweb/web/mime"
)

func TestInfer(t *testing.T) {
	assert.Equal(t, mime.HTML, mime.Infer("somepage.html"))
	assert.Equal(t, mime.HTML, mime.Infer("somepage.htm"))
	assert.Equal(t, mime.HTML, mime.Infer("somePage.HTM"))

	assert.Equal(t, mime.JavaScript, mime.Infer("my/path/script.js"))
	assert.Equal(t, mime.JavaScript, mime.Infer("my/other/script.js.gz"))

	assert.Equal(t, mime.SVG, mime.Infer("/img/logo.svg"))
}
