package server

import (
	"expvar"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

var (
	requestCount = expvar.NewInt("request_count")
)

// LoggingMiddleware logs the Method and Path that are part of the http request.
func LoggingMiddleware(logger *logrus.Logger) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			logger.Infof("Request received: %s %s", r.Method, r.URL.Path)
			next.ServeHTTP(w, r)
		})
	}
}

// MetricsMiddleware register some metrics per every service call.
func MetricsMiddleware(log *logrus.Logger) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			startTime := time.Now()

			next.ServeHTTP(w, r)

			duration := time.Since(startTime)

			requestCount.Add(1)

			log.Infof("Request(%d): %s %s took %v\n", requestCount.Value(), r.Method, r.URL.Path, duration)
		})
	}
}
