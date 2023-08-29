package handler

import (
	"fmt"
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
func ValidateUser(userService service.User, jwtService service.JWTManager, logger *logrus.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		qv := r.URL.Query()

		username := qv.Get("username")
		password := qv.Get("pwd")

		isValid, err := userService.Validate(username, password)
		if err != nil {
			logger.Errorf("user validation error: %s", err.Error())
			util.SendErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		if !isValid {
			msg := "user doesn't exists"
			logger.Error(msg)
			util.SendErrorResponse(w, http.StatusNotFound, msg)
			return
		}

		stringToken, err := jwtService.GenerateToken(username)
		if err != nil {
			fmt.Println(err)
			util.SendErrorResponse(w, http.StatusInternalServerError, err.Error())
		}

		resp := util.MapResponse{
			"status": "ok",
			"token":  "Bearer " + stringToken,
		}

		header := map[string]string{
			"Authorization": "Bearer " + stringToken,
		}

		respJSON := util.CustomMarshall(w, resp, logger)
		util.SendJSONResponseWithHeaders(w, http.StatusOK, respJSON, header)
	}
}
