package controllers

import (
	"golang-block-chain/services"
	"golang-block-chain/web/cores"
	"net/http"
)

func useDefaultErrorHandler(err error) (*cores.Response, error) {
	switch err := err.(type) {
	case services.ServiceError:
		var statusCode int
		switch err.(type) {
		case *services.ValidateError:
			statusCode = http.StatusBadRequest
		case *services.NotExistError:
			statusCode = http.StatusNotFound
		case *services.DuplicateError:
			statusCode = http.StatusConflict
		default:
			statusCode = http.StatusInternalServerError
		}
		resp := cores.NewErrorResponse(statusCode, err.Error(), err.Details())
		return resp, nil
	default:
		return nil, err
	}
}
