// Package service defines services and messages handled by the server.
package service

import (
	"encoding/json"
	"log"

	"github.com/golang/protobuf/proto"
	"github.com/toba/coreweb/socket"
)

type (
	// Service is a function that processes a client message. Each module
	// defines a set of services it provides.
	Service func(req *Request) *Response

	// Request contains the service call payload unmarshalled from the raw client
	// request.
	Request struct {
		Payload interface{}
	}

	// Endpoint defines how a service can be called -- its contract.
	Endpoint struct {
		Service              Service
		AllowAnonymous       bool
		RequireAuthorization []string
		Expect               interface{}
	}

	// ServiceMap matches endpoints to an identifier.
	ServiceMap map[ServiceID]*Endpoint

	// rawRequest sent by client to be processed by a module. Payload is
	// unmarshalled to retrieve service ID then matched to an endpoint that
	// supplies the PayloadType to unmarshal the raw payload into a specific
	// service request.
	rawRequest struct {
		ServiceID ServiceID `json:"type"`
		// RequestID is an identifier the client may use to match the response
		// back to a particular component action.
		RequestID string          `json:"id"`
		Payload   json.RawMessage `json:"data"`
	}

	// Response sent to the client. It is first marshalled to a JSON byte array.
	Response struct {
		StatusID  ServiceStatus `json:"status"`
		RequestID string        `json:"id"`
		Payload   interface{}   `json:"data"`
	}

	// ServiceID identifies the service method to be invoked.
	ServiceID int

	// ServiceStatus is returned with the response for every WebSocket request.
	ServiceStatus int
)

const (
	Okay ServiceStatus = iota
	IncompatiblePayload
	EmptyPayload
	UnableToParseRequest
	InvalidService
	UnableToMarshalResponse
	NotImplemented
	DatabaseError
	LdapError
	NoWebSocketClient
	NoMatchingRecords
)

// Handle processes an incoming WebSocket request by matching it to a service
// endpoint.
//
// The raw request is first unmarshalled to a standard service request that is
// used to lookup the service endpoint which provides the interface type to
// unmarshal the request payload.
//
// The payload is then unmarshalled to the indicated type and used to invoke
// the service function itself.
func Handle(endpoints ServiceMap) socket.RequestHandler {
	return func(socketRequest *socket.Request) []byte {
		var res *Response
		raw := &ServiceRequest{}
		err := proto.Unmarshal(socketRequest.Message, raw)

		if err != nil {
			res = Error(UnableToParseRequest)
		} else if ep, exists := endpoints[raw.ServiceID]; exists {
			if ep.Expect != nil {
				value := ep.Expect
				if err = json.Unmarshal(raw.Payload, value); err != nil {
					log.Print(err.Error())
					res = Error(IncompatiblePayload)
				} else if socketRequest.Client == nil {
					res = Error(NoWebSocketClient)
				} else {
					res = ep.Service(&Request{
						Payload: value,
					})
				}
			} else {
				res = ep.Service(&Request{})
			}
		} else {
			res = Error(InvalidService)
		}

		res.RequestID = raw.RequestID

		return res.JSON()
	}
}

// JSON converts the Response to a JSON byte array to be sent to the browser.
func (r *Response) JSON() []byte {
	data, err := json.Marshal(r)
	if err != nil {
		data, _ = json.Marshal(&Response{
			StatusID:  UnableToMarshalResponse,
			RequestID: r.RequestID,
			Payload:   nil,
		})
	}
	return data
}

// Error returns a response with an error status and no payload.
func Error(status ServiceStatus) *Response {
	return &Response{StatusID: status}
}

// Success returns a response with an okay status and the
// given payload.
func Success(payload interface{}) *Response {
	return &Response{StatusID: Okay, Payload: payload}
}

// MakeEndpoint defines the path, permissions and data shape expected by a
// service endpoint.
func MakeEndpoint(s Service, expect interface{}) *Endpoint {
	return &Endpoint{Service: s, Expect: expect}
}

// RequireLogin indicates if the endpoint requires authentication. The default
// is true.
func (ep *Endpoint) RequireLogin(require bool) *Endpoint {
	ep.AllowAnonymous = !require
	return ep
}
