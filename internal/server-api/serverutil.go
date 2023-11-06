package serverapi

import "net/http"

func onError(err error, w http.ResponseWriter, status int) bool {
	if err != nil {
		http.Error(w, err.Error(), status)
		return true
	}
	return false
}

func renderError(w http.ResponseWriter, message string, statusCode int) {
	w.WriteHeader(statusCode)
	w.Write([]byte(message))
}

func validateType(fileType string, w http.ResponseWriter) bool {
	switch fileType {
	case "image/jpeg", "image/jpg", "image/gif", "image/png", "application/pdf", "video/mp4":
		return true
	default:
		return false
	}
}
