package api

type (
	ResponseSuccess struct {
		StatusCode int64       `json:"status_code"`
		Message    string      `json:"message"`
		Data       interface{} `json:"data"`
	}

	ResponseError struct {
		StatusCode int64       `json:"status_code"`
		Message    string      `json:"message"`
		Error      interface{} `json:"error"`
	}
)
