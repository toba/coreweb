package coreweb_test

// https://elithrar.github.io/article/testing-http-handlers-go/
import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/toba/coreweb"
	"github.com/toba/coreweb/auth"
	"github.com/toba/coreweb/encoding"
	"github.com/toba/coreweb/header/accept"
	"github.com/toba/coreweb/header/content"
	"github.com/toba/coreweb/mime"
)

var (
	c = web.Config{
		FromFolder: "static",
	}
	modulePaths = []string{"module1", "module2"}
	authPaths   = map[string]*auth.AuthProvider{
		"auth/dropbox": auth.Providers[auth.Dropbox],
	}
	handler = web.Handle(c, modulePaths, authPaths)
)

func get(t *testing.T, path string) *http.Response {
	assert.NotNil(t, handler)

	r := httptest.NewRequest(http.MethodGet, path, nil)
	w := httptest.NewRecorder()
	r.Header.Add(accept.Encoding, encoding.GZip)

	handler(w, r)

	return w.Result()
}

func TestMissingFile(t *testing.T) {
	res := get(t, "/no-such-file")
	assert.Equal(t, http.StatusNotFound, res.StatusCode)
}

func TestModulePath(t *testing.T) {
	res := get(t, "/module1")
	assert.Equal(t, http.StatusOK, res.StatusCode)
	assert.Equal(t, mime.HTML, res.Header.Get(content.Type))
}

func TestLogo(t *testing.T) {
	res := get(t, "/img/logo.svg")
	assert.Equal(t, http.StatusOK, res.StatusCode)
	assert.Equal(t, mime.SVG, res.Header.Get(content.Type))
}

func TestJavascript(t *testing.T) {
	res := get(t, "/js/common.js")
	assert.Equal(t, http.StatusOK, res.StatusCode)
	assert.Equal(t, mime.JavaScript, res.Header.Get(content.Type))
	assert.Equal(t, encoding.GZip, res.Header.Get(content.Encoding))
}
