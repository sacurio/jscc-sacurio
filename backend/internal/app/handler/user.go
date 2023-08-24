package handler

import (
	"net/http"

	"github.com/sacurio/jb-challenge/internal/app/service"
	"github.com/sacurio/jb-challenge/internal/util"
	"github.com/sirupsen/logrus"
)

// DefaultHandler...
func DefaultHandler(logger *logrus.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data := util.MapResponse{
			"status": "ok",
		}

		dataJSON := util.CustomMarshall(w, data, logger)
		util.SendJSONResponse(w, http.StatusOK, dataJSON)
	}
}

// ValidateUser...
func ValidateUser(userService service.User, logger *logrus.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		qv := r.URL.Query()

		username := qv.Get("username")
		password := qv.Get("pwd")

		isValid, err := userService.Validate(username, password)
		if err != nil {
			logger.Errorf("user validation error: %s", err.Error())
			util.SendErrorResponse(w, http.StatusInternalServerError, err.Error())
		}

		if !isValid {
			msg := "user doesn't exists"
			logger.Error(msg)
			util.SendErrorResponse(w, http.StatusNotFound, msg)
			return
		}

		resp := util.MapResponse{
			"status": "ok",
		}

		respJSON := util.CustomMarshall(w, resp, logger)
		util.SendJSONResponse(w, http.StatusOK, respJSON)
	}
}
