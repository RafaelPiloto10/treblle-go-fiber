package treblle_fiber

import (
	"encoding/json"
	"net/http/httptest"
	"time"
)

type ResponseInfo struct {
	Headers  map[string]string      `json:"headers"`
	Code     int                    `json:"code"`
	Size     int                    `json:"size"`
	LoadTime float64                `json:"load_time"`
	Body     map[string]interface{} `json:"body"`
	Errors   []ErrorInfo            `json:"errors"`
}

type ErrorInfo struct {
	Source  string `json:"source"`
	Type    string `json:"type"`
	Message string `json:"message"`
	File    string `json:"file"`
	Line    int    `json:"line"`
}

// Extract information from the response recorder
func getFiberResponseInfo(response *httptest.ResponseRecorder, startTime time.Time) ResponseInfo {
	defer dontPanic()
	responseBytes := response.Body.Bytes()

	errInfo := ErrorInfo{}
	var body map[string]interface{}

	// handle error down below
	err := json.Unmarshal(responseBytes, &body)

	headers := make(map[string]string)
	for k := range response.Header() {
		headers[k] = response.Header().Get(k)
	}

	r := ResponseInfo{
		Headers:  headers,
		Code:     response.Code,
		Size:     len(responseBytes),
		LoadTime: float64(time.Since(startTime).Microseconds()),
		Body:     body,
		Errors:   make([]ErrorInfo, 0),
	}

	if err != nil {
		errInfo.Message = err.Error()
		r.Errors = append(r.Errors, errInfo)
	}

	return r
}
