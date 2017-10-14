package service_test

import (
	"math"
	"testing"

	"github.com/toba/coreweb/socket"
	"toba.tech/app/lib/module"

	"strconv"

	"encoding/json"

	"github.com/stretchr/testify/assert"
)

type testPayload struct {
	Number1 int
	Number2 int
	Text    string
}

const (
	Test1 module.ServiceID = iota
	Test2
	Test3
)

var text string

var services = module.ServiceMap{
	Test1: &module.Endpoint{
		AllowAnonymous: true,
		Expect:         &testPayload{},
		Service: func(req *module.Request) *module.Response {
			p := req.Payload.(*testPayload)
			return module.Success(p.Number1 + p.Number2)
		},
	},

	Test2: &module.Endpoint{
		AllowAnonymous: true,
		Expect:         &text,
		Service: func(req *module.Request) *module.Response {
			return module.Error(module.DatabaseError)
		},
	},
}

func respond(t *testing.T, payload string) *module.Response {
	handler := module.Handle(services)
	req := &socket.Request{
		Message: []byte(payload),
		Client:  &socket.Client{},
	}
	res := handler(req)
	assert.NotNil(t, res)

	rx := &module.Response{}
	err := json.Unmarshal(res, rx)
	assert.NoError(t, err)
	assert.NotNil(t, rx)
	return rx
}

func TestEndpointMatch(t *testing.T) {
	res := respond(t, `{
		"type": `+strconv.Itoa(int(Test1))+`,
		"id": "refID",
		"data": {
			"number1": 1,
			"number2": 2,
			"text": "random"
		}
	}`)
	assert.Equal(t, "refID", res.RequestID)
	assert.Equal(t, module.Okay, res.StatusID)
	assert.Equal(t, float64(3), res.Payload)

	res = respond(t, `{
		"type": `+strconv.Itoa(int(Test2))+`,
		"id": "refID",
		"data": "something"
	}`)
	assert.Equal(t, module.DatabaseError, res.StatusID)
}

func TestInvalidEndpoint(t *testing.T) {
	res := respond(t, `{
		"type": `+strconv.Itoa(int(Test3))+`,
		"id": "refID",
		"data": null
	}`)
	assert.Equal(t, module.InvalidService, res.StatusID)
}

func TestBadPayload(t *testing.T) {
	res := respond(t, `{
		"type": `+strconv.Itoa(int(Test1))+`,
		"id": "refID",
		"data": ""
	}`)
	assert.Equal(t, "refID", res.RequestID)
	assert.Equal(t, module.IncompatiblePayload, res.StatusID)
}

func TestJSON(t *testing.T) {
	res := &module.Response{
		RequestID: "23",
		StatusID:  module.Okay,
	}
	expect := []byte(`{"status":` + strconv.Itoa(int(module.Okay)) + `,"id":"23","data":null}`)
	json := res.JSON()

	assert.Equal(t, expect, json)

	// channels cannot be marshelled to JSON
	res = &module.Response{
		RequestID: "24",
		StatusID:  module.Okay,
		Payload:   make(chan int),
	}
	expect = []byte(`{"status":` + strconv.Itoa(int(module.UnableToMarshalResponse)) + `,"id":"24","data":null}`)
	json = res.JSON()

	assert.Equal(t, expect, json)

	// or infinity
	res = &module.Response{
		RequestID: "24",
		StatusID:  module.Okay,
		Payload:   math.Inf(1),
	}
	json = res.JSON()

	assert.Equal(t, expect, json)
}
