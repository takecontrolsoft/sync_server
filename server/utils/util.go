package utils

import (
	"math/rand"
	"net/http"
	"strings"
	"time"

	l "github.com/takecontrolsoft/logger"
)

func RenderIfError(err error, w http.ResponseWriter, statusCode int) bool {
	if err != nil {
		http.Error(w, err.Error(), statusCode)
		l.LogError(err)
		return true
	}
	return false
}

func RenderMessage(w http.ResponseWriter, message string, statusCode int) {
	w.WriteHeader(statusCode)
	w.Write([]byte(message))
	l.LogMessage(message)
}

func IsAllowedFileType(fileType string, w http.ResponseWriter) bool {
	allowed := []string{"image/", "video/", "audio/"}
	result := false
	for i := range allowed {
		result = result || strings.HasPrefix(fileType, allowed[i])
	}
	return result
}

func GenerateRandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	seed := rand.NewSource(time.Now().UnixNano())
	random := rand.New(seed)

	result := make([]byte, length)
	for i := range result {
		result[i] = charset[random.Intn(len(charset))]
	}
	return string(result)
}
