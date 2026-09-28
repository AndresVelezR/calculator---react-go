package httpapi

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestCORS(t *testing.T) {
	testCases := []struct {
		name            string
		method          string
		origin          string
		requestedMethod string
		expectedStatus  int
		expectedOrigin  string
	}{
		{"allowed calculation", "POST", "http://localhost:5173", "", 200, "http://localhost:5173"},
		{"no origin", "POST", "", "", 200, ""},
		{"other origin gets no permission", "POST", "http://example.com", "", 200, ""},
		{"preflight", "OPTIONS", "http://localhost:5173", "POST", 204, "http://localhost:5173"},
		{"GET preflight", "OPTIONS", "http://localhost:5173", "GET", 204, "http://localhost:5173"},
		{"unapproved preflight", "OPTIONS", "http://example.com", "POST", 403, ""},
		{"unapproved method", "OPTIONS", "http://localhost:5173", "DELETE", 405, "http://localhost:5173"},
		{"OPTIONS without origin", "OPTIONS", "", "POST", 405, ""},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			request := httptest.NewRequest(testCase.method, "/api/v1/calculate", strings.NewReader(`{"operation":"add","a":1,"b":2}`))
			request.Header.Set("Origin", testCase.origin)
			request.Header.Set("Access-Control-Request-Method", testCase.requestedMethod)
			request.Header.Set("Access-Control-Request-Headers", "content-type")
			recorder := httptest.NewRecorder()
			CORS(NewRouter(), "http://localhost:5173").ServeHTTP(recorder, request)
			if recorder.Code != testCase.expectedStatus {
				t.Errorf("status = %d, expected %d", recorder.Code, testCase.expectedStatus)
			}
			if origin := recorder.Header().Get("Access-Control-Allow-Origin"); origin != testCase.expectedOrigin {
				t.Errorf("allowed origin = %q, expected %q", origin, testCase.expectedOrigin)
			}
			if recorder.Header().Get("Vary") != "Origin" {
				t.Error("response must vary by Origin")
			}
			if testCase.expectedStatus == http.StatusNoContent {
				if recorder.Header().Get("Access-Control-Allow-Headers") != "Content-Type" ||
					recorder.Header().Get("Access-Control-Allow-Methods") != "GET, POST, OPTIONS" || recorder.Body.Len() != 0 {
					t.Error("preflight must permit JSON requests and return an empty body")
				}
			}
		})
	}
}

func TestCORSCustomOriginOnError(t *testing.T) {
	request := httptest.NewRequest(http.MethodPost, "/api/v1/calculate", strings.NewReader(`{"operation":"divide","a":1,"b":0}`))
	request.Header.Set("Origin", "http://localhost:3000")
	recorder := httptest.NewRecorder()
	CORS(NewRouter(), "http://localhost:3000").ServeHTTP(recorder, request)
	if recorder.Code != http.StatusUnprocessableEntity || recorder.Header().Get("Access-Control-Allow-Origin") != "http://localhost:3000" {
		t.Error("custom origin must also receive CORS headers on error responses")
	}
}
