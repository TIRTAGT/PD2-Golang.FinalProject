package handlerstruct

import "net/http"

type ControllerHandler func(w http.ResponseWriter, r *http.Request)(*string)

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