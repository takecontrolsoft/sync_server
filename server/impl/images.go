/* Copyright 2024 Take Control - Software & Infrastructure

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
	"encoding/json"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/disintegration/imaging"
	"github.com/go-errors/errors"
	"github.com/takecontrolsoft/sync_server/server/config"
	"github.com/takecontrolsoft/sync_server/server/utils"
)

type fileData struct {
	UserData userData
	File     string
	Quality  string
}

func GetImageHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		var result fileData
		if err := json.NewDecoder(r.Body).Decode(&result); err != nil {
			utils.RenderError(w, errors.Errorf("$Required json input {UserData: { User: '', DeviceId: ''}, 	File: ''}"), http.StatusBadRequest)
			return
		}
		userFromClient := result.UserData.User
		deviceId := result.UserData.DeviceId
		userId := ResolveToUserId(userFromClient)
		if userId == "" {
			userId = userFromClient
		}
		file := result.File
		quality := result.Quality
		userDirName := filepath.Join(config.UploadDirectory, userId, deviceId)
		originalFilePath := filepath.Join(userDirName, file)
		if quality == "full" {
			// Serve original file as-is — no decode/re-encode, no quality change.
			if err := serveOriginalFile(w, originalFilePath, file); err != nil {
				utils.RenderError(w, err, http.StatusInternalServerError)
			}
			return
		}

		path := ""
		thumbnailAddedExtension, err := utils.GetThumbnailFileAddedExtension(originalFilePath)
		if err != nil {
			utils.RenderError(w, err, http.StatusInternalServerError)
			return
		}
		path = fmt.Sprintf("%s%s", ThumbnailBasePath(userDirName, file), thumbnailAddedExtension)
		src, err := utils.GetImageFromFilePath(path)
		if err != nil {
			utils.RenderError(w, err, http.StatusInternalServerError)
			return
		}

		if quality == "high" {
			metadataPath := MetadataPath(userDirName, file)
			orientation := GetOrientationFromMetadata(metadataPath)
			src = applyEXIFOrientation(src, orientation)
			src = resizeMaxLongEdge(src, 1920)
			w.Header().Set("Content-Type", "image/jpeg")
			w.WriteHeader(http.StatusOK)
			jpeg.Encode(w, src, &jpeg.Options{Quality: 85})
			return
		}

		// Thumbnail: PNG
		w.Header().Set("Content-Type", "image/png")
		w.WriteHeader(http.StatusOK)
		png.Encode(w, src)
	}
}

// serveOriginalFile streams the file unchanged; Content-Type from extension.
func serveOriginalFile(w http.ResponseWriter, filePath, file string) error {
	f, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer f.Close()
	contentType := contentTypeFromFileName(file)
	w.Header().Set("Content-Type", contentType)
	w.WriteHeader(http.StatusOK)
	_, err = io.Copy(w, f)
	return err
}

func contentTypeFromFileName(path string) string {
	ext := strings.ToLower(filepath.Ext(path))
	switch ext {
	case ".jpg", ".jpeg":
		return "image/jpeg"
	case ".png":
		return "image/png"
	case ".gif":
		return "image/gif"
	case ".webp":
		return "image/webp"
	case ".heic":
		return "image/heic"
	default:
		return "application/octet-stream"
	}
}

// resizeMaxLongEdge resizes the image so the longest edge is at most maxPx; aspect ratio preserved.
func resizeMaxLongEdge(src image.Image, maxPx int) image.Image {
	b := src.Bounds()
	w, h := b.Dx(), b.Dy()
	if w <= maxPx && h <= maxPx {
		return src
	}
	if w >= h {
		return imaging.Resize(src, maxPx, 0, imaging.Lanczos)
	}
	return imaging.Resize(src, 0, maxPx, imaging.Lanczos)
}
