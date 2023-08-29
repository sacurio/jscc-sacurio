package server

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/sirupsen/logrus"
)

func TestLoggingMiddleware(t *testing.T) {
	mockLogger := logrus.New()
	mockLogger.Out = nil

	var logBuffer bytes.Buffer
	mockLogger.SetOutput(&logBuffer)

	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Test-Header", "test")
	})

	middleware := LoggingMiddleware(mockLogger)(handler)

	middleware.ServeHTTP(rr, req)

	expectedLogEntry := "Request: GET /"
	if strings.Contains(logBuffer.String(), expectedLogEntry) {
		t.Errorf("Expected logger entry: %s, but got: %s", expectedLogEntry, logBuffer.String())
	}

	expectedHeaderValue := "test"
	if rr.Header().Get("X-Test-Header") != expectedHeaderValue {
		t.Errorf("Expected response header value: %s, but got: %s", expectedHeaderValue, rr.Header().Get("X-Test-Header"))
	}
}

func TestAuthMiddleware(t *testing.T) {
	mockLogger := logrus.New()
	mockLogger.Out = nil // Suppress logger output during testing

	var logBuffer bytes.Buffer
	mockLogger.SetOutput(&logBuffer)

	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Dummy handler logic
		w.WriteHeader(http.StatusOK)
	})

	middleware := AuthMiddleware(mockLogger)(handler)

	req.Header.Set("Authorization", "Bearer valid_token")

	middleware.ServeHTTP(rr, req)

	expectedLogMessage := "Getting JWT from client..."
	if !strings.Contains(logBuffer.String(), expectedLogMessage) {
		t.Errorf("Expected log message: %s, but got: %s", expectedLogMessage, logBuffer.String())
	}

	req.Header.Del("Authorization")

	logBuffer.Reset()

	middleware.ServeHTTP(rr, req)

	if rr.Code != http.StatusUnauthorized {
		t.Errorf("Expected response status code: %d, but got: %d", http.StatusUnauthorized, rr.Code)
	}

	expectedErrorMessage := "user auth is not valid"
	if strings.Contains(logBuffer.String(), expectedErrorMessage) {
		t.Errorf("Expected error message: %s, but got: %s", expectedErrorMessage, logBuffer.String())
	}
}

func TestEnableCORSMiddleware(t *testing.T) {
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Dummy handler logic
		w.WriteHeader(http.StatusOK)
	})

	middleware := EnableCORSMiddleware(handler)

	middleware.ServeHTTP(rr, req)

	expectedAllowOrigin := "*"
	if rr.Header().Get("Access-Control-Allow-Origin") != expectedAllowOrigin {
		t.Errorf("Expected Access-Control-Allow-Origin header: %s, but got: %s",
			expectedAllowOrigin, rr.Header().Get("Access-Control-Allow-Origin"))
	}

	expectedAllowMethods := "GET, POST, PUT, DELETE, OPTIONS"
	if rr.Header().Get("Access-Control-Allow-Methods") != expectedAllowMethods {
		t.Errorf("Expected Access-Control-Allow-Methods header: %s, but got: %s",
			expectedAllowMethods, rr.Header().Get("Access-Control-Allow-Methods"))
	}

	expectedAllowHeaders := "Content-Type, Authorization"
	if rr.Header().Get("Access-Control-Allow-Headers") != expectedAllowHeaders {
		t.Errorf("Expected Access-Control-Allow-Headers header: %s, but got: %s",
			expectedAllowHeaders, rr.Header().Get("Access-Control-Allow-Headers"))
	}

	if rr.Code != http.StatusOK {
		t.Errorf("Expected response status code: %d, but got: %d", http.StatusOK, rr.Code)
	}
}
