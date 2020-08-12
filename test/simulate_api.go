package test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

// SimulationData - used to minify API tests
type SimulationData struct {
	Method             string
	API                string
	Body               []byte
	ExpectedStatusCode int
	ErrorMessage       string
	Cookie             http.Cookie
}

// SetRequestHeaderJSON - as name suggests
func SetRequestHeaderJSON(req *http.Request) {
	req.Header.Set("Content-Type", "application/json")
}

// ResponseError - Error message if not same status code as expected
func ResponseError(t *testing.T, message string, expected, got int, body *bytes.Buffer) {
	t.Errorf("%v\n\nExpected status code %v, got %v\nBody: %v", message, expected, got, body)
}

// SimulateAPI - simulate API tests
func SimulateAPI(t *testing.T, router *gin.Engine, s SimulationData) *httptest.ResponseRecorder {
	w := httptest.NewRecorder()
	http.SetCookie(w, &s.Cookie)
	jsonBody := bytes.NewBuffer(s.Body)
	req, _ := http.NewRequest(s.Method, s.API, jsonBody)
	req.Header = http.Header{"Cookie": w.HeaderMap["Set-Cookie"]}
	SetRequestHeaderJSON(req)
	router.ServeHTTP(w, req)

	if w.Code != s.ExpectedStatusCode {
		ResponseError(t, s.ErrorMessage, s.ExpectedStatusCode, w.Code, w.Body)
	}

	return w
}
