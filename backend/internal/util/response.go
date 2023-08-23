package util

import (
	"encoding/json"
	"net/http"

	"github.com/sirupsen/logrus"
)

// MapResponse provides a way to send structured data to client through a pair key-value.
type MapResponse map[string]interface{}

// SendErrorResponse writes an error value to ResponseWriter and setting the StatusInternalServerError status code.
func SendErrorResponse(w http.ResponseWriter, statusCode int, err string) {
	dataErr := MapResponse{
		"error": err,
	}

	errJSON, errMsg := json.Marshal(dataErr)
	if errMsg != nil {
		sendResponse(w, http.StatusInternalServerError, []byte(errMsg.Error()))
	}

	sendResponse(w, statusCode, errJSON)
}

// SendJSONResponse writes data to ResponseWriter and setting the http status code received as well.
func SendJSONResponse(w http.ResponseWriter, statusCode int, data []byte) {
	sendResponse(w, statusCode, data)
}

func sendResponse(w http.ResponseWriter, statusCode int, data []byte) {
	w.WriteHeader(statusCode)
	w.Write(data)
}

// CustomMarshall send data received in JSON format, logging if error exists.
func CustomMarshall(w http.ResponseWriter, data any, logger *logrus.Logger) []byte {
	respJSON, err := json.Marshal(data)
	if err != nil {
		logger.Error(err.Error())
		SendErrorResponse(w, http.StatusInternalServerError, err.Error())
		return nil
	}

	return respJSON
}
