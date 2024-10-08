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

package impl

import (
	"bufio"
	"encoding/json"
	"io"
	"mime"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/flytam/filenamify"
	"github.com/takecontrolsoft/go_multi_log/logger"
	"github.com/takecontrolsoft/sync_server/server/config"
	"github.com/takecontrolsoft/sync_server/server/mediatypes"
	"github.com/takecontrolsoft/sync_server/server/utils"
)

// Upload file handler for uploading large streamed files.
// A new file is saved under a directory named like the client device.
// An error will be rendered in the response if:
// - the file already exists;
// - the maximum allowed size is exceeded;
// - the file format is not allowed;
func UploadHandler(w http.ResponseWriter, r *http.Request) {
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
	mediatype, err := validateFileType(b, w)
	if utils.RenderIfError(err, w, http.StatusBadRequest) {
		return
	}

	userNameEncoded := r.Header.Get("user")

	var name []byte
	if err := json.Unmarshal([]byte(userNameEncoded), &name); err != nil {
		utils.RenderError(w, err, http.StatusInternalServerError)
		return
	}
	userName := string(name)
	if len(userName) == 0 {
		utils.RenderError(w, MissingUser, http.StatusBadRequest)
		return
	}
	dateClassifier := r.Header.Get("date")
	if len(dateClassifier) == 0 {
		utils.RenderError(w, MissingDateClassifier, http.StatusBadRequest)
		return
	}

	dateArray := strings.Split(dateClassifier, "-")
	if len(dateArray) < 2 {
		utils.RenderError(w, WrongDateClassifier, http.StatusBadRequest)
		return
	}
	year := dateArray[0]
	month := dateArray[1]

	_, params, err := mime.ParseMediaType(mp.Header.Get("Content-Disposition"))
	if err != nil {
		utils.RenderError(w, WrongDateClassifier, http.StatusBadRequest)
		return
	}

	deviceId, err := filenamify.Filenamify(params["name"], filenamify.Options{})
	if err != nil {
		utils.RenderError(w, WrongDateClassifier, http.StatusBadRequest)
		return
	}

	filename, err := filenamify.Filenamify(params["filename"], filenamify.Options{})
	if err != nil {
		utils.RenderError(w, WrongDateClassifier, http.StatusBadRequest)
		return
	}

	f, err := createNewFile(mp, w, userName, filename, deviceId, year, month)
	if utils.RenderIfError(err, w, http.StatusInternalServerError) {
		return
	}
	defer f.Close()

	var maxSize int64 = config.MaxUploadFileSize
	lmt := io.MultiReader(b, io.LimitReader(mp, maxSize-511))
	written, err := io.Copy(f, lmt)
	if utils.RenderIfError(err, w, http.StatusInternalServerError) {
		return
	}
	if written > maxSize {
		os.Remove(f.Name())
		utils.RenderError(w, FileSizeExceeded, http.StatusBadRequest)
		return
	}
	var relPath = filepath.Join(year, month, filename)

	go func() {
		_, err = ExtractMetadata(userName, deviceId, relPath)
		if err != nil {
			logger.ErrorF("Creating metadata failed for file %s, %v", relPath, err)
		}
		switch mediatype {
		case mediatypes.Video:
			_, err = BuildVideoThumbnail(userName, deviceId, relPath)
		case mediatypes.Image:
			_, err = BuildImageThumbnail(userName, deviceId, relPath)
		case mediatypes.Audio:
			_, err = BuildAudioThumbnail(userName, deviceId, relPath)
		default:
			logger.Info("Unknown media type for thumbnail")
		}
		if err != nil {
			logger.ErrorF("Creating thumbnail failed for file %s, %v", relPath, err)
		}
	}()

}

func createNewFile(mp *multipart.Part, w http.ResponseWriter,
	userName string, filename string, deviceId string,
	year string, month string) (*os.File, error) {

	dirName := filepath.Join(config.UploadDirectory, userName, deviceId, year, month)
	err := os.MkdirAll(dirName, os.ModePerm)
	if err != nil {
		return nil, err
	}
	filePath := filepath.Join(dirName, filename)
	f, err := os.Create(filePath)
	if err != nil {
		return nil, err
	}
	return f, err
}

func validateFileType(b *bufio.Reader, w http.ResponseWriter) (mediatypes.MediaType, error) {
	n, _ := b.Peek(512)
	fileType := http.DetectContentType(n)
	mediaType := utils.GetMediaType(fileType)
	if !utils.IsAllowedFileType(fileType, w) {
		err := InvalidFileTypeUploaded(fileType)
		return mediaType, err
	}
	return mediaType, nil
}
