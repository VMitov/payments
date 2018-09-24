package main

import (
	"github.com/VMitov/payments/pkg/errors"
	"github.com/go-chi/render"
)

func errRender(err error) render.Renderer {
	return &errors.ErrResponse{
		Err:            err,
		HTTPStatusCode: 422,
		StatusText:     "Error rendering response.",
		ErrorText:      err.Error(),
	}
}
