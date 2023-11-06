package serverapi

import (
	"bufio"
	"io"
	"mime"
	"net/http"
	"os"
	"path/filepath"
)

// Upload file handler for uploading large streamed files.
func uploadFileHandler() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.Body = http.MaxBytesReader(w, r.Body, maxUploadFileSize)
		reader, err := r.MultipartReader()
		if onError(err, w, http.StatusBadRequest) {
			return
		}
		mp, err := reader.NextPart()
		if onError(err, w, http.StatusInternalServerError) {
			return
		}

		_, params, err := mime.ParseMediaType(mp.Header.Get("Content-Disposition"))
		if onError(err, w, http.StatusInternalServerError) {
			return
		}
		deviceId := params["name"]
		filename := params["filename"]

		b := bufio.NewReader(mp)
		n, _ := b.Peek(512)
		fileType := http.DetectContentType(n)
		if !validateType(fileType, w) {
			renderError(w, "INVALID_FILE_TYPE", http.StatusBadRequest)
			return
		}

		dirName := filepath.Join(uploadDirectory, deviceId)
		err = os.MkdirAll(dirName, os.ModePerm)
		if onError(err, w, http.StatusInternalServerError) {
			return
		}
		filePath := filepath.Join(dirName, filename)
		f, err := os.Create(filePath)
		if onError(err, w, http.StatusInternalServerError) {
			return
		}
		defer f.Close()
		var maxSize int64 = maxUploadFileSize
		lmt := io.MultiReader(b, io.LimitReader(mp, maxSize-511))
		written, err := io.Copy(f, lmt)
		if onError(err, w, http.StatusInternalServerError) {
			return
		}
		if written > maxSize {
			os.Remove(f.Name())
			renderError(w, "FILE_SIZE_EXCEEDED", http.StatusBadRequest)
			return
		}
	})
}
