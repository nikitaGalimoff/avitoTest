package handler

import (
	"net/http"

	"avitotest/internal/domain"

	"github.com/labstack/echo/v4"
)

type ErrorResponse struct {
	Error struct {
		Code    string `json:"code"`
		Message string `json:"message"`
	} `json:"error"`
}

func WriteError(c echo.Context, err error, statusCode int) error {
	domainErr, ok := err.(*domain.DomainError)
	if !ok {
		if err != nil {
			domainErr = domain.NewDomainError(domain.ErrorCodeNotFound, err.Error())
		} else {
			domainErr = domain.NewDomainError(domain.ErrorCodeNotFound, "unknown error")
		}
		if statusCode == 0 {
			statusCode = http.StatusInternalServerError
		}
	}

	if statusCode == 0 {
		switch domainErr.Code {
		case domain.ErrorCodeNotFound:
			statusCode = http.StatusNotFound
		case domain.ErrorCodeTeamExists, domain.ErrorCodePRExists:
			statusCode = http.StatusConflict
		case domain.ErrorCodePRMerged, domain.ErrorCodeNotAssigned, domain.ErrorCodeNoCandidate:
			statusCode = http.StatusConflict
		default:
			statusCode = http.StatusBadRequest
		}
	}

	resp := ErrorResponse{
		Error: struct {
			Code    string `json:"code"`
			Message string `json:"message"`
		}{
			Code:    string(domainErr.Code),
			Message: domainErr.Message,
		},
	}

	return c.JSON(statusCode, resp)
}

func WriteJSON(c echo.Context, statusCode int, data interface{}) error {
	return c.JSON(statusCode, data)
}
