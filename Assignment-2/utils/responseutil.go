package utils

import (
	"encoding/json"
	"io"
	"net/http"
)

type BaseResponse struct {
	Code    int               `json:"code"`
	Message string            `json:"message"`
	Errors  ErrorResponseData `json:"error"`
	Data    interface{}       `json:"data"`
}

func (br *BaseResponse) ToJSON(w io.Writer) error {
	return json.NewEncoder(w).Encode(br)
}

func (br *BaseResponse) SendResponse(rw *http.ResponseWriter) error {
	(*rw).WriteHeader(br.Code)
	return json.NewEncoder(*rw).Encode(br)
}

// untuk nyimpen respon yg dkirim dlm bentuk json
func NewErrorResponseValue(code string, value string) ErrorData {
	return ErrorData{
		Code:    code,
		Message: value,
	}
}

func NewErrorResponseData(errorResponsesValue ...ErrorData) ErrorResponseData {
	errors := ErrorResponseData{}

	for _, v := range errorResponsesValue {
		errors = append(errors, v)
	}
	return errors
}

func NewErrorResponse(code int, message string, errors ...ErrorData) *BaseResponse {
	return &BaseResponse{
		Code:    code,
		Message: message,
		Errors: NewErrorResponseData(
			errors...,
		),
		Data: nil,
	}
}

func NewBaseResponse(code int, message string, errors ErrorResponseData, data interface{}) *BaseResponse {
	return &BaseResponse{
		Code:    code,
		Message: message,
		Errors:  errors,
		Data:    data,
	}
}

type ErrorData struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type ErrorResponseData []ErrorData

func (er *ErrorResponseData) ToJSON(w io.Writer) error {
	return json.NewEncoder(w).Encode(er)
}
