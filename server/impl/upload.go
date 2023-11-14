/* Copyright 2023 Take Control - Software & Infrastructure

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Package impl providers the main implementation of sync server APIs.
package impl

import (
	"bufio"
	"io"
	"mime"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"

	"github.com/takecontrolsoft/sync/server/config"
	"github.com/takecontrolsoft/sync/server/utils"
)

// Upload file handler for uploading large streamed files.
// A new file is saved under a directory named like the client device.
// An error will be rendered in the response if:
// - the file already exists;
// - the maximum allowed size is exceeded;
// - the file format is not allowed;
func UploadHandler() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.Body = http.MaxBytesReader(w, r.Body, config.MaxUploadFileSize)
		reader, err := r.MultipartReader()
		if utils.RenderIfError(err, w, http.StatusBadRequest) {
			return
		}
		mp, err := reader.NextPart()
		if utils.RenderIfError(err, w, http.StatusInternalServerError) {
			return
		}

		b := bufio.NewReader(mp)
		success := validateFileType(b, w)
		if !success {
			return
		}
		f, success := createNewFile(mp, w)
		if !success {
			return
		}
		var maxSize int64 = config.MaxUploadFileSize
		lmt := io.MultiReader(b, io.LimitReader(mp, maxSize-511))
		written, err := io.Copy(f, lmt)
		if utils.RenderIfError(err, w, http.StatusInternalServerError) {
			return
		}
		if written > maxSize {
			os.Remove(f.Name())
			utils.RenderMessage(w, "FILE_SIZE_EXCEEDED", http.StatusBadRequest)
			return
		}
	})
}

func createNewFile(mp *multipart.Part, w http.ResponseWriter) (*os.File, bool) {
	_, params, err := mime.ParseMediaType(mp.Header.Get("Content-Disposition"))
	if utils.RenderIfError(err, w, http.StatusInternalServerError) {
		return nil, false
	}
	deviceId := params["name"]
	filename := params["filename"]

	dirName := filepath.Join(config.UploadDirectory, deviceId)
	err = os.MkdirAll(dirName, os.ModePerm)
	if utils.RenderIfError(err, w, http.StatusInternalServerError) {
		return nil, false
	}
	filePath := filepath.Join(dirName, filename)
	f, err := os.Create(filePath)
	if utils.RenderIfError(err, w, http.StatusInternalServerError) {
		return nil, false
	}
	defer f.Close()
	return f, true
}

func validateFileType(b *bufio.Reader, w http.ResponseWriter) bool {
	n, _ := b.Peek(512)
	fileType := http.DetectContentType(n)
	if !utils.IsAllowedFileType(fileType, w) {
		utils.RenderMessage(w, "INVALID_FILE_TYPE", http.StatusBadRequest)
		return false
	}
	return true
}
