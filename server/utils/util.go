package utils

import (
	"math/rand"
	"net/http"
	"strings"
	"time"

	"github.com/takecontrolsoft/go_multi_log/logger"
)

func RenderIfError(err error, w http.ResponseWriter, statusCode int) bool {
	if err != nil {
		http.Error(w, err.Error(), statusCode)
		logger.Error(err)
		return true
	}
	return false
}

func RenderError(w http.ResponseWriter, err error, statusCode int) {
	w.WriteHeader(statusCode)
	w.Write([]byte(err.Error()))
	logger.Error(err)
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
