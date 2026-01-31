/* Copyright 2026 Take Control - Software & Infrastructure

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
	"fmt"
	"io"
	"mime"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/flytam/filenamify"
	"github.com/takecontrolsoft/go_multi_log/logger"
	"github.com/takecontrolsoft/sync_server/server/auth"
	"github.com/takecontrolsoft/sync_server/server/config"
	"github.com/takecontrolsoft/sync_server/server/media"
	"github.com/takecontrolsoft/sync_server/server/mediatypes"
	"github.com/takecontrolsoft/sync_server/server/paths"
	"github.com/takecontrolsoft/sync_server/server/trash"
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
	userFromClient := string(name)
	if len(userFromClient) == 0 {
		utils.RenderError(w, MissingUser, http.StatusBadRequest)
		return
	}
	userId := auth.ResolveUserId(userFromClient)
	if userId == "" {
		userId = userFromClient
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
	// Clamp future or bogus year/month to current date so files don't end up in future-year folders
	year, month = clampYearMonth(year, month)

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

	saveToTrash := strings.EqualFold(strings.TrimSpace(r.Header.Get("X-Save-To-Trash")), "true")

	f, err := createNewFile(mp, w, userId, filename, deviceId, year, month, saveToTrash)
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
	var relPath string
	if saveToTrash {
		relPath = filepath.Join(paths.TrashFolder, year, month, filename)
	} else {
		relPath = filepath.Join(year, month, filename)
	}
	// Use forward slashes so ThumbnailBasePath/MetadataPath recognize Trash paths on all OSes.
	relPath = paths.Normalize(relPath)

	go func() {
		userDir := filepath.Join(config.UploadDirectory, userId, deviceId)
		// 1. Wait for metadata creation to complete.
		_, metaErr := media.ExtractMetadata(userId, deviceId, relPath)
		if metaErr != nil {
			logger.ErrorF("Creating metadata failed for file %s, %v", relPath, metaErr)
		}
		// 2. Wait for thumbnail creation to complete.
		var thumbErr error
		switch mediatype {
		case mediatypes.Video:
			_, thumbErr = media.BuildVideoThumbnail(userId, deviceId, relPath)
		case mediatypes.Image:
			_, thumbErr = media.BuildImageThumbnail(userId, deviceId, relPath)
		case mediatypes.Audio:
			_, thumbErr = media.BuildAudioThumbnail(userId, deviceId, relPath)
		default:
			logger.Info("Unknown media type for thumbnail")
		}
		if thumbErr != nil {
			logger.ErrorF("Creating thumbnail failed for file %s, %v", relPath, thumbErr)
		}
		// 3. Run document-to-trash detection only after both metadata and thumbnail have completed,
		// so the file, metadata, and thumbnail are all moved to Trash/ together.
		if metaErr == nil && thumbErr == nil && mediatype == mediatypes.Image && config.DocumentToTrashEnabled && !saveToTrash {
			fullPath := filepath.Join(userDir, relPath)
			if config.DocumentClassifierPath != "" {
				RunDocumentClassifierSync(fullPath, userDir, relPath)
			} else if LooksLikeDocument(fullPath) {
				trash.MoveToTrash(userDir, relPath)
			}
		}
	}()

}

func createNewFile(mp *multipart.Part, w http.ResponseWriter,
	userId string, filename string, deviceId string,
	year string, month string, saveToTrash bool) (*os.File, error) {

	dirName := filepath.Join(config.UploadDirectory, userId, deviceId, year, month)
	if saveToTrash {
		dirName = filepath.Join(config.UploadDirectory, userId, deviceId, paths.TrashFolder, year, month)
	}
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

// clampYearMonth returns (year, month) strings; if year is in the future or before 2000,
// returns current year and month so uploads don't create future-year or bogus folders.
func clampYearMonth(year, month string) (string, string) {
	yr, errY := strconv.Atoi(year)
	mn, errM := strconv.Atoi(month)
	now := time.Now()
	if errY != nil || errM != nil || yr > now.Year() || yr < 2000 || mn < 1 || mn > 12 {
		return fmt.Sprintf("%d", now.Year()), fmt.Sprintf("%d", int(now.Month()))
	}
	return year, month
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
