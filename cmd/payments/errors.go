package main

import (
	"net/http"

	"github.com/VMitov/payments/pkg/errors"
	"github.com/go-chi/render"
)

var errNotFound = &errors.ErrResponse{
	HTTPStatusCode: http.StatusNotFound,
	StatusText:     "Resource not found.",
}

func errInvalidRequest(err error) render.Renderer {
	return &errors.ErrResponse{
		Err:            err,
		HTTPStatusCode: http.StatusBadRequest,
		StatusText:     "Invalid request.",
		ErrorText:      err.Error(),
	}
}

func errSystem(err error) render.Renderer {
	return &errors.ErrResponse{
		Err:            err,
		HTTPStatusCode: http.StatusInternalServerError,
		StatusText:     "Error handling response.",
		ErrorText:      err.Error(),
	}
}
