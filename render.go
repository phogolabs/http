package rest

import (
	"net/http"

	"github.com/apex/log"
	"github.com/go-chi/render"
	"github.com/phogolabs/rest/middleware"
)

func init() {
	render.Decode = Decode
	render.Respond = Respond
}

// Binder interface for managing request payloads.
type Binder interface {
	Bind(r *http.Request) error
}

// Renderer interface for managing response payloads.
type Renderer interface {
	Render(w http.ResponseWriter, r *http.Request) error
}

// Bind decodes a request body and executes the Binder method of the
// payload structure.
func Bind(r *http.Request, v Binder) error {
	return render.Bind(r, v)
}

// Render renders a single payload and respond to the client request.
func Render(w http.ResponseWriter, r *http.Request, v Renderer) error {
	return render.Render(w, r, v)
}

// Status sets a HTTP response status code hint into request context at any point
// during the request life-cycle. Before the Responder sends its response header
// it will check the StatusCtxKey
func Status(r *http.Request, status int) {
	render.Status(r, status)
}

// GetLogger returns the associated request logger
func GetLogger(r *http.Request) log.Interface {
	return middleware.GetLogger(r)
}