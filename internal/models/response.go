package models

type FailResponse struct {
	Error   string `json:"error"`
	Details string `json:"details,omitempty"`
}

type SuccesResponse struct {
	SuccesData interface{}
}

type ErrorResponse struct {
	Error error
}
