package socket_test

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/toba/goweb/web/socket"
	"toba.tech/app/lib/config"

	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
)

var (
	c        = config.HTTP{}
	hello    = []byte("hello")
	world    = []byte("world")
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
)

// https://play.golang.org/p/X8GLU-Gcox
func connect(t *testing.T, h socket.RequestHandler) *websocket.Conn {
	handler := socket.Handle(c, h)
	srv := httptest.NewServer(http.HandlerFunc(handler))
	u, _ := url.Parse(srv.URL)
	u.Scheme = "ws"

	conn, res, err := websocket.DefaultDialer.Dial(u.String(), nil)

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.NotNil(t, conn)

	assert.Equal(t, http.StatusSwitchingProtocols, res.StatusCode)
	assert.Contains(t, res.Header, socket.Accept)

	return conn
}

func mockHandler(t *testing.T) socket.RequestHandler {
	return func(req *socket.Request) []byte {
		assert.NotNil(t, req)
		assert.Equal(t, hello, req.Message)
		return world
	}
}

func TestServiceMessage(t *testing.T) {
	conn := connect(t, mockHandler(t))

	defer conn.Close()

	err := conn.WriteMessage(websocket.TextMessage, hello)
	assert.NoError(t, err)

	messageType, res, err := conn.ReadMessage()
	assert.NoError(t, err)
	assert.Equal(t, websocket.TextMessage, messageType)
	assert.Equal(t, world, res)
}
