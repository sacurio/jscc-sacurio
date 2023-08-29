package server

import (
	"expvar"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go/v4"
	"github.com/gorilla/mux"
	"github.com/sacurio/jb-challenge/internal/app/service"
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

func JWTMiddleware(jwtService service.JWTManager, logger *logrus.Logger) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			logger.Info("Handling JWT...")
			tokenString := r.Header.Get("Authorization")
			if tokenString == "" {
				logger.Warn("JWT is empty.")
			}

			token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
				sk, err := jwtService.SecretKey()
				if err != nil {
					logger.Errorf("token parsing error, %s", err.Error())
					// return nil, err
				}
				logger.Info("Secret Key generation successfully")
				return sk, nil
			})

			if err != nil || !token.Valid {
				logger.Warn("JWT is not valid.")
			}

			next.ServeHTTP(w, r)
		})
	}
}

func AuthMiddleware(logger *logrus.Logger) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			logger.Info("Getting JWT from client...")
			user := r.Header.Get("Authorization")
			if user != "" {
				logger.Error("user auth is not valid")
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func EnableCORSMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}
