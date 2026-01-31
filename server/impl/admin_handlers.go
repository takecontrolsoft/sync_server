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
	"encoding/json"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/takecontrolsoft/go_multi_log/logger"
	"github.com/takecontrolsoft/sync_server/server/auth"
	"github.com/takecontrolsoft/sync_server/server/config"
	"github.com/takecontrolsoft/sync_server/server/media"
	"github.com/takecontrolsoft/sync_server/server/paths"
	"github.com/takecontrolsoft/sync_server/server/trash"
	"github.com/takecontrolsoft/sync_server/server/utils"
)

// ListAllRelativeFiles returns relative paths (forward slashes) of all files under userDir,
// excluding Trash, Thumbnails, and Metadata directories.
func ListAllRelativeFiles(userDir string) ([]string, error) {
	var files []string
	err := filepath.WalkDir(userDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			rel, _ := filepath.Rel(userDir, path)
			rel = paths.Normalize(rel)
			// Skip walking into Trash, Thumbnails, Metadata
			if rel == paths.TrashFolder || rel == paths.ThumbnailsFolder || rel == paths.MetadataFolder ||
				strings.HasPrefix(rel, paths.TrashFolder+"/") || strings.HasPrefix(rel, paths.ThumbnailsFolder+"/") || strings.HasPrefix(rel, paths.MetadataFolder+"/") {
				return filepath.SkipDir
			}
			return nil
		}
		rel, err := filepath.Rel(userDir, path)
		if err != nil {
			return nil
		}
		files = append(files, paths.Normalize(rel))
		return nil
	})
	return files, err
}

// RegenerateThumbnailsHandler regenerates thumbnails for all media files (excluding Trash).
// POST body: { "UserData": { "User": "", "DeviceId": "" } }
func RegenerateThumbnailsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	var result struct {
		UserData userData `json:"UserData"`
	}
	if err := json.NewDecoder(r.Body).Decode(&result); err != nil {
		utils.RenderError(w, err, http.StatusBadRequest)
		return
	}
	userFromClient := result.UserData.User
	deviceId := result.UserData.DeviceId
	if userFromClient == "" || deviceId == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	userId := auth.ResolveUserId(userFromClient)
	if userId == "" {
		userId = userFromClient
	}
	userDir := filepath.Join(config.UploadDirectory, userId, deviceId)
	all, err := ListAllRelativeFiles(userDir)
	if err != nil {
		utils.RenderError(w, err, http.StatusInternalServerError)
		return
	}
	var regenerated int
	for _, rel := range all {
		if strings.HasPrefix(rel, paths.TrashFolder+"/") || rel == paths.TrashFolder {
			continue
		}
		ext := strings.ToLower(filepath.Ext(rel))
		if paths.IsImagePath(rel) {
			if _, err := media.BuildImageThumbnail(userId, deviceId, rel); err != nil {
				logger.ErrorF("Regenerate thumbnail %s: %v", rel, err)
			} else {
				regenerated++
			}
		} else if ext == ".mp4" || ext == ".mov" || ext == ".avi" || ext == ".mkv" || ext == ".webm" {
			if _, err := media.BuildVideoThumbnail(userId, deviceId, rel); err != nil {
				logger.ErrorF("Regenerate video thumbnail %s: %v", rel, err)
			} else {
				regenerated++
			}
		}
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]int{"Regenerated": regenerated})
}

// CleanOrphanThumbnailsHandler deletes thumbnail and metadata files that have no corresponding source file.
// POST body: { "UserData": { "User": "", "DeviceId": "" } }
func CleanOrphanThumbnailsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	var result struct {
		UserData userData `json:"UserData"`
	}
	if err := json.NewDecoder(r.Body).Decode(&result); err != nil {
		utils.RenderError(w, err, http.StatusBadRequest)
		return
	}
	userFromClient := result.UserData.User
	deviceId := result.UserData.DeviceId
	if userFromClient == "" || deviceId == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	userId := auth.ResolveUserId(userFromClient)
	if userId == "" {
		userId = userFromClient
	}
	userDir := filepath.Join(config.UploadDirectory, userId, deviceId)
	removed := cleanOrphanThumbnailsInDir(userDir, paths.ThumbnailsFolder)
	removed += cleanOrphanThumbnailsInDir(userDir, paths.TrashFolder+"/"+paths.ThumbnailsFolder)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]int{"Removed": removed})
}

// cleanOrphanThumbnailsInDir walks userDir/thumbSubdir (e.g. Thumbnails or Trash/Thumbnails) and removes
// thumbnail files whose source file no longer exists. Also removes corresponding metadata.
func cleanOrphanThumbnailsInDir(userDir, thumbSubdir string) int {
	dir := filepath.Join(userDir, filepath.FromSlash(thumbSubdir))
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		return 0
	}
	prefix := paths.Normalize(thumbSubdir) + "/"
	var removed int
	_ = filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return err
		}
		rel, err := filepath.Rel(userDir, path)
		if err != nil {
			return nil
		}
		rel = paths.Normalize(rel)
		if !strings.HasPrefix(rel, prefix) {
			return nil
		}
		sourceRel := strings.TrimPrefix(rel, prefix)
		if strings.HasSuffix(sourceRel, ".jpeg") {
			sourceRel = strings.TrimSuffix(sourceRel, ".jpeg")
		}
		var sourcePath string
		if strings.HasPrefix(thumbSubdir, paths.TrashFolder) {
			sourcePath = filepath.Join(userDir, paths.TrashFolder, sourceRel)
		} else {
			sourcePath = filepath.Join(userDir, sourceRel)
		}
		if _, err := os.Stat(sourcePath); err == nil {
			return nil
		}
		if err := os.Remove(path); err != nil {
			logger.ErrorF("Clean orphan thumbnail remove %s: %v", path, err)
			return nil
		}
		removed++
		var metaPath string
		if strings.HasPrefix(thumbSubdir, paths.TrashFolder) {
			metaPath = paths.MetadataPath(userDir, paths.TrashFolder+"/"+sourceRel)
		} else {
			metaPath = paths.MetadataPath(userDir, sourceRel)
		}
		_ = os.Remove(metaPath)
		return nil
	})
	return removed
}

// RunDocumentDetectionHandler runs document detection (Python classifier or built-in heuristic) on existing image files,
// moves detected documents to Trash (with thumbnails and metadata). Returns { "Moved": N }.
// POST body: { "UserData": { "User": "", "DeviceId": "" } }
func RunDocumentDetectionHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	var result struct {
		UserData userData `json:"UserData"`
	}
	if err := json.NewDecoder(r.Body).Decode(&result); err != nil {
		utils.RenderError(w, err, http.StatusBadRequest)
		return
	}
	userFromClient := result.UserData.User
	deviceId := result.UserData.DeviceId
	if userFromClient == "" || deviceId == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	userId := auth.ResolveUserId(userFromClient)
	if userId == "" {
		userId = userFromClient
	}
	userDir := filepath.Join(config.UploadDirectory, userId, deviceId)
	if _, err := os.Stat(userDir); os.IsNotExist(err) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(map[string]int{"Moved": 0})
		return
	}
	all, err := ListAllRelativeFiles(userDir)
	if err != nil {
		utils.RenderError(w, err, http.StatusInternalServerError)
		return
	}
	var moved int
	for _, rel := range all {
		if strings.HasPrefix(rel, paths.TrashFolder+"/") || rel == paths.TrashFolder {
			continue
		}
		if !paths.IsImagePath(rel) {
			continue
		}
		fullPath := filepath.Join(userDir, rel)
		classifierMoved := false
		if config.DocumentClassifierPath != "" {
			classifierMoved = RunDocumentClassifierSyncReturnsMoved(fullPath, userDir, rel)
		}
		if classifierMoved {
			moved++
		} else if LooksLikeDocument(fullPath) {
			trash.MoveToTrash(userDir, rel)
			moved++
		}
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]int{"Moved": moved})
}
