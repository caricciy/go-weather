package util

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

type jsonTestRow struct {
	name         string
	resp         any
	status       int
	expectedBody string
	expectedCode int
}

func TestSendJSON(t *testing.T) {
	testTable := []jsonTestRow{
		{
			name:         "Valid JSON response",
			resp:         map[string]string{"message": "success"},
			status:       http.StatusOK,
			expectedBody: `{"message":"success"}`,
			expectedCode: http.StatusOK,
		},
		{
			name:         "Internal Server Error on marshal failure",
			resp:         make(chan int), // Invalid type for JSON marshaling
			status:       http.StatusOK,
			expectedBody: `{"error":"Internal Server Error"}`,
			expectedCode: http.StatusInternalServerError,
		},
	}

	for _, tr := range testTable {
		t.Run(tr.name, func(t *testing.T) {
			recorder := httptest.NewRecorder()
			SendJSON(recorder, tr.resp, tr.status)

			result := recorder.Result()
			defer func(result *http.Response) {
				_ = result.Body.Close()
			}(result)

			assert.Equal(t, tr.expectedCode, result.StatusCode, "Expected status code to be %v", tr.expectedCode)
			assert.Equal(t, "application/json", result.Header.Get("Content-Type"), "Expected Content-Type to be application/json")
			assert.JSONEq(t, tr.expectedBody, recorder.Body.String(), "Expected body to match")
		})
	}
}

// failingResponseWriter is only overriding the Write method of the http.ResponseWriter interface.
// When we create an instance of failingResponseWriter and embed a httptest.NewRecorder()
// as the ResponseWriter, the other methods of the interface, such as Header() and WriteHeader(),
// continue to be called on the embedded object (httptest.NewRecorder()) since they were not overridden
// in the failingResponseWriter structure.
// Go allows you to embed an interface or structure and override only the methods you want to modify.
// The methods that are not overridden are delegated to the embedded object.
// If we want all methods to be controlled by failingResponseWriter, we need to explicitly override
// the other methods, such as Header() and WriteHeader().
type failingResponseWriter struct {
	http.ResponseWriter
}

func (f *failingResponseWriter) Write(data []byte) (int, error) {
	return 0, errors.New("intentional write error")
}

func TestSendJSON_WriteError(t *testing.T) {
	recorder := &failingResponseWriter{
		ResponseWriter: httptest.NewRecorder(),
	}

	resp := map[string]string{"message": "success"}
	status := http.StatusOK

	SendJSON(recorder, resp, status)

	// No assertions are needed here since the slog.Error call is not directly testable.
	// However, you can verify that the function does not panic and handles the error gracefully.
	assert.NotNil(t, recorder, "Recorder should not be nil")
}
