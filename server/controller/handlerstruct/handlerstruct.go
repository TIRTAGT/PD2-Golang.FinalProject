package handlerstruct

import "net/http"

type ControllerHandler func(w http.ResponseWriter, r *http.Request)

type ControllerStruct struct {
	// GET handler (optional)
	GET ControllerHandler

	// POST handler (optional)
	POST ControllerHandler

	// PUT handler (optional)
	PUT ControllerHandler

	// DELETE handler (optional)
	DELETE ControllerHandler
}