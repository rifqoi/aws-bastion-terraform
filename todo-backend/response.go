package main

import (
	"net/http"

	"github.com/go-chi/render"
	"github.com/rifqoi/aws-project/todo-backend/common"
	"github.com/rifqoi/aws-project/todo-backend/common/logs"
)

type SuccessResponse struct {
	Message    string `json:"message,omitempty"`
	StatusCode int    `json:"status_code,omitempty"`
	Data       any    `json:"data,omitempty"`
}

func (re SuccessResponse) Render(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(re.StatusCode)
	return nil
}

func NewSuccessResponse(msg string, data any) SuccessResponse {
	return SuccessResponse{
		Message: msg,
		Data:    data,
	}
}

// /////////////////////////////////
func httpErrorResponse(err string, errMsg string, code int, w http.ResponseWriter, r *http.Request) {
	resp := ErrorResponse{err, errMsg, code}

	render.Render(w, r, resp)
}

func InternalError(errMsg string, w http.ResponseWriter, r *http.Request) {
	httpErrorResponse(errMsg, "INTERNAL_SERVER_ERROR", http.StatusInternalServerError, w, r)
}

func BadRequest(errMsg string, w http.ResponseWriter, r *http.Request) {
	httpErrorResponse(errMsg, "BAD_REQUEST", http.StatusBadRequest, w, r)
}

func NotFound(errMsg string, w http.ResponseWriter, r *http.Request) {
	httpErrorResponse(errMsg, "NOT_FOUND", http.StatusNotFound, w, r)
}

func RespondHTTPError(err error, errMsg string, w http.ResponseWriter, r *http.Request) {
	slugErr, ok := err.(*common.Error)
	if !ok {
		InternalError(errMsg, w, r)
		return
	}

	// This is the convenient way of deciding which error should we sent.
	switch slugErr.Type() {
	case common.ErrorTypeInvalidInput:
		BadRequest(slugErr.Error(), w, r)
	case common.ErrorTypeNotFound:
		NotFound(errMsg, w, r)
	case common.ErrorTypeServerError:
		InternalError(errMsg, w, r)
	default:
		InternalError(errMsg, w, r)
	}

	logs.GetLogger().Error(slugErr)
}

type ErrorResponse struct {
	Error      string `json:"error,omitempty"`
	Message    string `json:"message,omitempty"`
	StatusCode int    `json:"-"`
}

func (e ErrorResponse) Render(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(e.StatusCode)
	return nil
}
