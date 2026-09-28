package httpapi

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestCalculateHandler(t *testing.T) {
	testCases := []struct {
		name           string
		body           string
		expectedStatus int
		expectedResult float64
		expectedError  string
	}{
		{"add", `{"operation":"add","a":1,"b":2}`, 200, 3, ""},
		{"subtract", `{"operation":"subtract","a":1,"b":2}`, 200, -1, ""},
		{"multiply", `{"operation":"multiply","a":3,"b":2}`, 200, 6, ""},
		{"divide", `{"operation":"divide","a":1,"b":2}`, 200, 0.5, ""},
		{"power", `{"operation":"power","a":2,"b":3}`, 200, 8, ""},
		{"square root without b", `{"operation":"square_root","a":9}`, 200, 3, ""},
		{"square root ignores b", `{"operation":"square_root","a":9,"b":-1}`, 200, 3, ""},
		{"percentage", `{"operation":"percentage","a":50,"b":20}`, 200, 10, ""},
		{"explicit zeroes", `{"operation":"add","a":0,"b":0}`, 200, 0, ""},
		{"empty body", ``, 400, 0, "request must be a valid JSON object with operation, a, and b"},
		{"malformed JSON", `{"operation":`, 400, 0, "request must be a valid JSON object with operation, a, and b"},
		{"number as string", `{"operation":"add","a":"1","b":2}`, 400, 0, "request must be a valid JSON object with operation, a, and b"},
		{"array", `[]`, 400, 0, "request must be a valid JSON object with operation, a, and b"},
		{"unknown field", `{"operation":"add","a":1,"b":2,"extra":3}`, 400, 0, "request must be a valid JSON object with operation, a, and b"},
		{"number outside float64", `{"operation":"add","a":1e999,"b":2}`, 400, 0, "request must be a valid JSON object with operation, a, and b"},
		{"oversized body", strings.Repeat(" ", 4097) + `{}`, 400, 0, "request must be a valid JSON object with operation, a, and b"},
		{"multiple objects", `{} {}`, 400, 0, "request must contain a single JSON object"},
		{"trailing garbage", `{} invalid`, 400, 0, "request must contain a single JSON object"},
		{"missing operation", `{"a":1,"b":2}`, 400, 0, "operation is required"},
		{"empty operation", `{"operation":"","a":1,"b":2}`, 400, 0, "operation is required"},
		{"null body", `null`, 400, 0, "operation is required"},
		{"unknown operation", `{"operation":"modulo","a":1,"b":2}`, 400, 0, "unknown operation"},
		{"missing a", `{"operation":"add","b":2}`, 400, 0, "a is required"},
		{"null a", `{"operation":"add","a":null,"b":2}`, 400, 0, "a is required"},
		{"missing b", `{"operation":"add","a":1}`, 400, 0, "b is required for this operation"},
		{"null b", `{"operation":"add","a":1,"b":null}`, 400, 0, "b is required for this operation"},
		{"divide by zero", `{"operation":"divide","a":1,"b":0}`, 422, 0, "division by zero"},
		{"negative square root", `{"operation":"square_root","a":-1}`, 422, 0, "square root of a negative number"},
		{"infinite result", `{"operation":"power","a":10,"b":1000}`, 422, 0, "result is not a finite number"},
		{"NaN result", `{"operation":"power","a":-1,"b":0.5}`, 422, 0, "result is not a finite number"},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodPost, "/api/v1/calculate", strings.NewReader(testCase.body))
			recorder := httptest.NewRecorder()
			NewRouter().ServeHTTP(recorder, request)
			if recorder.Code != testCase.expectedStatus {
				t.Fatalf("status = %d, expected %d; body = %s", recorder.Code, testCase.expectedStatus, recorder.Body.String())
			}
			if contentType := recorder.Header().Get("Content-Type"); contentType != "application/json" {
				t.Errorf("Content-Type = %q, expected application/json", contentType)
			}
			var response struct {
				Result *float64 `json:"result"`
				Error  string   `json:"error"`
			}
			if err := json.Unmarshal(recorder.Body.Bytes(), &response); err != nil {
				t.Fatalf("invalid JSON response: %v", err)
			}
			if response.Error != testCase.expectedError {
				t.Errorf("error = %q, expected %q", response.Error, testCase.expectedError)
			}
			if testCase.expectedStatus == http.StatusOK {
				if response.Result == nil || *response.Result != testCase.expectedResult {
					t.Errorf("body = %s, expected result %v", recorder.Body.String(), testCase.expectedResult)
				}
			} else if response.Result != nil {
				t.Error("error response must not contain a result")
			}
		})
	}
}

func TestHealth(t *testing.T) {
	recorder := httptest.NewRecorder()
	NewRouter().ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, "/health", nil))
	if recorder.Code != http.StatusOK || recorder.Body.String() != `{"status":"ok"}` {
		t.Errorf("health returned %d %s", recorder.Code, recorder.Body.String())
	}
}

func TestCalculateRejectsGet(t *testing.T) {
	recorder := httptest.NewRecorder()
	NewRouter().ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, "/api/v1/calculate", nil))
	if recorder.Code != http.StatusMethodNotAllowed {
		t.Errorf("status = %d, expected 405", recorder.Code)
	}
}
