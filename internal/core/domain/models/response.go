package models

type Token struct {
	JWT string `json:"jwt"`
}

type Response struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   *ErrorData  `json:"error,omitempty"`
}

type ErrorData struct {
	Message string            `json:"message"`
	Details string            `json:"details,omitempty"`
	Errors  []ValidationError `json:"errors,omitempty"`
}

type ValidationError struct {
	Field string `json:"field"`
	Tag   string `json:"tag"`
	Value string `json:"value"`
}
