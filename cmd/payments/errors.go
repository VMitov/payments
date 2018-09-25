package main

import (
	"github.com/VMitov/payments/pkg/errors"
	"github.com/go-chi/render"
)

var errNotFound = &errors.ErrResponse{
	HTTPStatusCode: 404,
	StatusText:     "Resource not found.",
}

func errInvalidRequest(err error) render.Renderer {
	return &errors.ErrResponse{
		Err:            err,
		HTTPStatusCode: 400,
		StatusText:     "Invalid request.",
		ErrorText:      err.Error(),
	}
}

func errRender(err error) render.Renderer {
	return &errors.ErrResponse{
		Err:            err,
		HTTPStatusCode: 422,
		StatusText:     "Error rendering response.",
		ErrorText:      err.Error(),
	}
}
