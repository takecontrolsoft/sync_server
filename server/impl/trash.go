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
	"runtime"
	"sort"
	"strings"

	"github.com/takecontrolsoft/sync_server/server/config"
	"github.com/takecontrolsoft/sync_server/server/utils"
)

const TrashFolder = "Trash"

// trashPrefix is "Trash/" for path logic (API uses forward slashes).
const trashPrefix = TrashFolder + "/"

// ThumbnailBasePath returns the full path to the thumbnail file (without .jpeg extension for videos).
// For files in Trash, thumbnails live under Trash/Thumbnails/; otherwise under Thumbnails/.
func ThumbnailBasePath(userDir, file string) string {
	if strings.HasPrefix(file, trashPrefix) {
		rest := strings.TrimPrefix(file, trashPrefix)
		return filepath.Join(userDir, TrashFolder, "Thumbnails", rest)
	}
	return filepath.Join(userDir, "Thumbnails", file)
}

// MetadataPath returns the full path to the metadata JSON file.
// For files in Trash, metadata lives under Trash/Metadata/; otherwise under Metadata/.
func MetadataPath(userDir, file string) string {
	if strings.HasPrefix(file, trashPrefix) {
		rest := strings.TrimPrefix(file, trashPrefix)
		return filepath.Join(userDir, TrashFolder, "Metadata", rest+".json")
	}
	return filepath.Join(userDir, "Metadata", file+".json")
}

type moveToTrashData struct {
	UserData userData
	Files    []string
}

type restoreData struct {
	UserData userData
	Files    []string
}

// MoveToTrashHandler moves files (and their thumbnails and metadata) to Trash.
// POST body: { "UserData": { "User": "", "DeviceId": "" }, "Files": ["2024/01/photo.jpg", ...] }
func MoveToTrashHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	var result moveToTrashData
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
	userId := ResolveToUserId(userFromClient)
	if userId == "" {
		userId = userFromClient
	}
	userDir := filepath.Join(config.UploadDirectory, userId, deviceId)

	for _, file := range result.Files {
		if file == "" || strings.Contains(file, "..") {
			continue
		}
		// Skip if already in Trash (API uses forward slash: "Trash/...")
		if strings.HasPrefix(file, trashPrefix) || file == TrashFolder {
			continue
		}
		originalPath := filepath.Join(userDir, file)
		trashPath := filepath.Join(userDir, TrashFolder, file)

		// Get thumbnail extension while main file still exists (videos use .jpeg).
		thumbExt, _ := utils.GetThumbnailFileAddedExtension(originalPath)
		thumbSrc := filepath.Join(userDir, "Thumbnails", file) + thumbExt
		thumbDst := filepath.Join(userDir, TrashFolder, "Thumbnails", file) + thumbExt

		// Move main file first; then thumbnail and metadata (no-op if src missing).
		if err := moveFile(originalPath, trashPath); err != nil {
			continue
		}
		_ = moveFile(thumbSrc, thumbDst)
		metaSrc := filepath.Join(userDir, "Metadata", file+".json")
		metaDst := filepath.Join(userDir, TrashFolder, "Metadata", file+".json")
		_ = moveFile(metaSrc, metaDst)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

// MoveRelativePathToTrash moves one file (and its thumbnail and metadata) from
// the normal folder to Trash. relPath is e.g. "2024/01/photo.jpg". No-op if
// relPath is already under Trash. Used by upload when document detection is enabled.
func MoveRelativePathToTrash(userDir, relPath string) {
	if relPath == "" || strings.Contains(relPath, "..") {
		return
	}
	// Skip if already in Trash; accept both "Trash/" and OS separator
	if strings.HasPrefix(relPath, trashPrefix) || relPath == TrashFolder {
		return
	}
	if strings.HasPrefix(relPath, TrashFolder+string(os.PathSeparator)) {
		return
	}
	originalPath := filepath.Join(userDir, relPath)
	trashPath := filepath.Join(userDir, TrashFolder, relPath)
	thumbExt, _ := utils.GetThumbnailFileAddedExtension(originalPath)
	thumbSrc := filepath.Join(userDir, "Thumbnails", relPath) + thumbExt
	thumbDst := filepath.Join(userDir, TrashFolder, "Thumbnails", relPath) + thumbExt
	if moveFile(originalPath, trashPath) != nil {
		return
	}
	_ = moveFile(thumbSrc, thumbDst)
	metaSrc := filepath.Join(userDir, "Metadata", relPath+".json")
	metaDst := filepath.Join(userDir, TrashFolder, "Metadata", relPath+".json")
	_ = moveFile(metaSrc, metaDst)
}

// RestoreHandler moves files from Trash back to their original folder (by path).
// POST body: { "UserData": { "User": "", "DeviceId": "" }, "Files": ["Trash/2024/01/photo.jpg", ...] }
func RestoreHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	var result restoreData
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
	userId := ResolveToUserId(userFromClient)
	if userId == "" {
		userId = userFromClient
	}
	userDir := filepath.Join(config.UploadDirectory, userId, deviceId)
	// Use package trashPrefix "Trash/" — API always sends forward slashes.

	for _, file := range result.Files {
		if file == "" || strings.Contains(file, "..") {
			continue
		}
		if !strings.HasPrefix(file, trashPrefix) {
			continue
		}
		restorePath := strings.TrimPrefix(file, trashPrefix)
		trashFilePath := filepath.Join(userDir, TrashFolder, restorePath)
		originalFilePath := filepath.Join(userDir, restorePath)

		// Move main file back
		if err := moveFile(trashFilePath, originalFilePath); err != nil {
			continue
		}

		// Move thumbnail back from Trash/Thumbnails
		thumbExt, _ := utils.GetThumbnailFileAddedExtension(originalFilePath)
		thumbTrash := filepath.Join(userDir, TrashFolder, "Thumbnails", restorePath) + thumbExt
		thumbOriginal := filepath.Join(userDir, "Thumbnails", restorePath) + thumbExt
		_ = moveFile(thumbTrash, thumbOriginal)

		// Move metadata back from Trash/Metadata
		metaTrash := filepath.Join(userDir, TrashFolder, "Metadata", restorePath+".json")
		metaOriginal := filepath.Join(userDir, "Metadata", restorePath+".json")
		_ = moveFile(metaTrash, metaOriginal)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

type emptyTrashData struct {
	User     string `json:"User"`
	DeviceId string `json:"DeviceId"`
	Password string `json:"Password"`
}

// EmptyTrashHandler permanently deletes the entire Trash folder (and all its contents:
// main files, Trash/Thumbnails, Trash/Metadata). Requires auth (token or password).
// POST body: { "User": "", "DeviceId": "", "Password": "" } or Authorization: Bearer <token>
func EmptyTrashHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	var result emptyTrashData
	if err := json.NewDecoder(r.Body).Decode(&result); err != nil {
		utils.RenderError(w, err, http.StatusBadRequest)
		return
	}
	userFromClient := result.User
	deviceId := result.DeviceId
	if userFromClient == "" || deviceId == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	userId := RequireAuthForDangerous(w, r, userFromClient, deviceId, result.Password)
	if userId == "" {
		return
	}
	userDir := filepath.Join(config.UploadDirectory, userId, deviceId)
	trashPath := filepath.Join(userDir, TrashFolder)
	if err := removeTrashDirAll(trashPath); err != nil {
		utils.RenderError(w, err, http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

// removeTrashDirAll deletes the entire Trash directory and its contents.
// On Windows, clears read-only before delete so os.RemoveAll-style deletion succeeds.
func removeTrashDirAll(trashPath string) error {
	if _, err := os.Stat(trashPath); err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	var files []string
	var dirs []string
	err := filepath.WalkDir(trashPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			dirs = append(dirs, path)
			return nil
		}
		files = append(files, path)
		return nil
	})
	if err != nil {
		return err
	}
	// Remove all files first (on Windows clear read-only so delete succeeds).
	for _, path := range files {
		if runtime.GOOS == "windows" {
			_ = os.Chmod(path, 0666)
		}
		if err := os.Remove(path); err != nil {
			return err
		}
	}
	// Remove subdirs from deepest to shallowest, then the Trash root.
	sort.Slice(dirs, func(i, j int) bool { return len(dirs[i]) > len(dirs[j]) })
	for _, path := range dirs {
		if runtime.GOOS == "windows" {
			_ = os.Chmod(path, 0777)
		}
		if err := os.Remove(path); err != nil {
			return err
		}
	}
	return nil
}

// moveFile moves a file, creating parent dirs of dst. No-op if src does not exist.
func moveFile(src, dst string) error {
	if _, err := os.Stat(src); os.IsNotExist(err) {
		return nil
	}
	if err := os.MkdirAll(filepath.Dir(dst), 0755); err != nil {
		return err
	}
	return os.Rename(src, dst)
}

// ListTrashFiles returns relative paths of main files under userDir/Trash (e.g. "Trash/2024/01/photo.jpg").
// Skips Trash/Thumbnails and Trash/Metadata so only media files are listed.
func ListTrashFiles(userDir string) ([]string, error) {
	var files []string
	trashDir := filepath.Join(userDir, TrashFolder)
	if _, err := os.Stat(trashDir); os.IsNotExist(err) {
		return files, nil
	}
	err := filepath.WalkDir(trashDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		rel, err := filepath.Rel(userDir, path)
		if err != nil {
			return nil
		}
		rel = filepath.ToSlash(rel)
		// Only list main files; skip Trash/Thumbnails/... and Trash/Metadata/...
		if strings.HasPrefix(rel, TrashFolder+"/Thumbnails/") || strings.HasPrefix(rel, TrashFolder+"/Metadata/") {
			return nil
		}
		files = append(files, rel)
		return nil
	})
	return files, err
}
