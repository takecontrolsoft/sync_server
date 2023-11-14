package utils

import (
	"internal/errors_util"
	"net/http"
)

func RenderIfError(err error, w http.ResponseWriter, statusCode int) bool {
	if err != nil {
		http.Error(w, err.Error(), statusCode)
		errors_util.LogError(err)
		return true
	}
	return false
}

func RenderMessage(w http.ResponseWriter, message string, statusCode int) {
	w.WriteHeader(statusCode)
	w.Write([]byte(message))
	errors_util.LogMessage(message)
}

func IsAllowedFileType(fileType string, w http.ResponseWriter) bool {
	switch fileType {
	case "image/jpeg", "image/jpg", "image/gif", "image/png", "application/pdf", "video/mp4":
		return true
	default:
		return false
	}
}
