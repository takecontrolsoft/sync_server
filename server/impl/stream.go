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
	"bufio"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/takecontrolsoft/sync_server/server/config"
	"github.com/takecontrolsoft/sync_server/server/utils"
)

// GetStreamHandler serves raw media files (video/audio) with HTTP Range support
// for streaming and seeking without full download.
// GET /stream?User=...&DeviceId=...&File=... (File is URL-encoded path, e.g. 2024/01/video.mp4)
func GetStreamHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	userFromClient := r.URL.Query().Get("User")
	deviceId := r.URL.Query().Get("DeviceId")
	file := r.URL.Query().Get("File")
	if userFromClient == "" || deviceId == "" || file == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	// Prevent path traversal: file must not contain ".."
	if strings.Contains(file, "..") {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	userId := ResolveToUserId(userFromClient)
	if userId == "" {
		userId = userFromClient
	}
	// Normalize path: server expects forward slashes (Windows clients may send backslash).
	file = strings.ReplaceAll(file, "\\", "/")
	userDirName := filepath.Join(config.UploadDirectory, userId, deviceId)
	originalFilePath := filepath.Join(userDirName, file)
	// Ensure resolved path is still under userDirName
	absPath, err := filepath.Abs(originalFilePath)
	if err != nil {
		utils.RenderError(w, err, http.StatusInternalServerError)
		return
	}
	absUserDir, _ := filepath.Abs(userDirName)
	if !strings.HasPrefix(absPath, absUserDir) {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	f, err := os.Open(originalFilePath)
	if err != nil {
		if os.IsNotExist(err) {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		utils.RenderError(w, err, http.StatusInternalServerError)
		return
	}
	defer f.Close()
	info, err := f.Stat()
	if err != nil {
		utils.RenderError(w, err, http.StatusInternalServerError)
		return
	}
	if info.IsDir() {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	size := info.Size()
	contentType := getContentType(originalFilePath, f)
	_, _ = f.Seek(0, 0) // reset after getContentType may have read from f
	w.Header().Set("Content-Type", contentType)
	w.Header().Set("Accept-Ranges", "bytes")
	w.Header().Set("Content-Disposition", "inline") // play in place, do not download
	rangeHeader := r.Header.Get("Range")
	if rangeHeader == "" {
		w.Header().Set("Content-Length", strconv.FormatInt(size, 10))
		w.WriteHeader(http.StatusOK)
		_, _ = io.Copy(w, f)
		return
	}
	// Parse "bytes=start-end"
	if !strings.HasPrefix(rangeHeader, "bytes=") {
		w.WriteHeader(http.StatusRequestedRangeNotSatisfiable)
		return
	}
	rangeStr := strings.TrimPrefix(rangeHeader, "bytes=")
	parts := strings.Split(rangeStr, "-")
	var start, end int64
	if len(parts) == 1 {
		start, err = strconv.ParseInt(strings.TrimSpace(parts[0]), 10, 64)
		if err != nil {
			w.WriteHeader(http.StatusRequestedRangeNotSatisfiable)
			return
		}
		end = size - 1
	} else {
		start, err = strconv.ParseInt(strings.TrimSpace(parts[0]), 10, 64)
		if err != nil {
			start = 0
		}
		if parts[1] == "" {
			end = size - 1
		} else {
			end, err = strconv.ParseInt(strings.TrimSpace(parts[1]), 10, 64)
			if err != nil || end >= size {
				end = size - 1
			}
		}
	}
	if start > end || start < 0 {
		w.WriteHeader(http.StatusRequestedRangeNotSatisfiable)
		return
	}
	if start >= size {
		w.WriteHeader(http.StatusRequestedRangeNotSatisfiable)
		return
	}
	if end >= size {
		end = size - 1
	}
	contentLength := end - start + 1
	_, err = f.Seek(start, 0)
	if err != nil {
		utils.RenderError(w, err, http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Length", strconv.FormatInt(contentLength, 10))
	w.Header().Set("Content-Range", "bytes "+strconv.FormatInt(start, 10)+"-"+strconv.FormatInt(end, 10)+"/"+strconv.FormatInt(size, 10))
	w.WriteHeader(http.StatusPartialContent)
	_, _ = io.CopyN(w, f, contentLength)
}

func getContentType(filePath string, f *os.File) string {
	// Prefer detection from content for accuracy
	buf := bufio.NewReader(f)
	peek, _ := buf.Peek(512)
	ct := http.DetectContentType(peek)
	// Only use for media; avoid "application/octet-stream" for known extensions
	if strings.HasPrefix(ct, "video/") || strings.HasPrefix(ct, "audio/") {
		return ct
	}
	ext := strings.ToLower(filepath.Ext(filePath))
	switch ext {
	case ".mp4", ".m4v":
		return "video/mp4"
	case ".webm":
		return "video/webm"
	case ".mov":
		return "video/quicktime"
	case ".3gp":
		return "video/3gpp"
	case ".mkv":
		return "video/x-matroska"
	case ".mp3":
		return "audio/mpeg"
	case ".m4a":
		return "audio/mp4"
	case ".ogg", ".oga":
		return "audio/ogg"
	case ".wav":
		return "audio/wav"
	default:
		if ct != "application/octet-stream" {
			return ct
		}
		return "application/octet-stream"
	}
}
