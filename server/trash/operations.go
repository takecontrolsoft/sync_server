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

// Package trash provides operations for moving files to and from the trash folder.
package trash

import (
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/takecontrolsoft/sync_server/server/paths"
	"github.com/takecontrolsoft/sync_server/server/utils"
)

const trashPrefix = paths.TrashFolder + "/"

// MoveToTrash moves a file and its associated thumbnail/metadata to Trash.
// userDir is the base directory (UploadDirectory/userId/deviceId).
// relPath is the relative path within userDir (e.g., "2024/01/photo.jpg").
func MoveToTrash(userDir, relPath string) {
	if relPath == "" || strings.Contains(relPath, "..") {
		return
	}
	if strings.HasPrefix(relPath, trashPrefix) || relPath == paths.TrashFolder {
		return
	}
	if strings.HasPrefix(relPath, paths.TrashFolder+string(os.PathSeparator)) {
		return
	}
	originalPath := filepath.Join(userDir, relPath)
	trashPath := filepath.Join(userDir, paths.TrashFolder, relPath)
	thumbExt, _ := utils.GetThumbnailFileAddedExtension(originalPath)
	thumbSrc := filepath.Join(userDir, paths.ThumbnailsFolder, relPath) + thumbExt
	thumbDst := filepath.Join(userDir, paths.TrashFolder, paths.ThumbnailsFolder, relPath) + thumbExt
	if moveFile(originalPath, trashPath) != nil {
		return
	}
	_ = moveFile(thumbSrc, thumbDst)
	metaSrc := filepath.Join(userDir, paths.MetadataFolder, relPath+".json")
	metaDst := filepath.Join(userDir, paths.TrashFolder, paths.MetadataFolder, relPath+".json")
	_ = moveFile(metaSrc, metaDst)
}

// MoveFileToTrash moves a single file (without thumbnail/metadata) to Trash.
// Returns error if move fails.
func MoveFileToTrash(userDir, relPath string) error {
	if relPath == "" || strings.Contains(relPath, "..") {
		return nil
	}
	if strings.HasPrefix(relPath, trashPrefix) || relPath == paths.TrashFolder {
		return nil
	}
	originalPath := filepath.Join(userDir, relPath)
	trashPath := filepath.Join(userDir, paths.TrashFolder, relPath)
	return moveFile(originalPath, trashPath)
}

// MoveWithAssociatedFiles moves a file and its thumbnail/metadata to Trash.
// Used by HTTP handlers. Returns error if the main file move fails.
func MoveWithAssociatedFiles(userDir, relPath string) error {
	if relPath == "" || strings.Contains(relPath, "..") {
		return nil
	}
	if strings.HasPrefix(relPath, trashPrefix) || relPath == paths.TrashFolder {
		return nil
	}
	originalPath := filepath.Join(userDir, relPath)
	trashPath := filepath.Join(userDir, paths.TrashFolder, relPath)

	thumbExt, _ := utils.GetThumbnailFileAddedExtension(originalPath)
	thumbSrc := filepath.Join(userDir, paths.ThumbnailsFolder, relPath) + thumbExt
	thumbDst := filepath.Join(userDir, paths.TrashFolder, paths.ThumbnailsFolder, relPath) + thumbExt

	if err := moveFile(originalPath, trashPath); err != nil {
		return err
	}
	_ = moveFile(thumbSrc, thumbDst)
	metaSrc := filepath.Join(userDir, paths.MetadataFolder, relPath+".json")
	metaDst := filepath.Join(userDir, paths.TrashFolder, paths.MetadataFolder, relPath+".json")
	_ = moveFile(metaSrc, metaDst)
	return nil
}

// RestoreFromTrash restores a file and its thumbnail/metadata from Trash.
// trashRelPath should NOT include the "Trash/" prefix.
// Returns error if the main file restore fails.
func RestoreFromTrash(userDir, trashRelPath string) error {
	if trashRelPath == "" || strings.Contains(trashRelPath, "..") {
		return nil
	}
	trashFilePath := filepath.Join(userDir, paths.TrashFolder, trashRelPath)
	originalFilePath := filepath.Join(userDir, trashRelPath)

	if err := moveFile(trashFilePath, originalFilePath); err != nil {
		return err
	}

	thumbExt, _ := utils.GetThumbnailFileAddedExtension(originalFilePath)
	thumbTrash := filepath.Join(userDir, paths.TrashFolder, paths.ThumbnailsFolder, trashRelPath) + thumbExt
	thumbOriginal := filepath.Join(userDir, paths.ThumbnailsFolder, trashRelPath) + thumbExt
	_ = moveFile(thumbTrash, thumbOriginal)

	metaTrash := filepath.Join(userDir, paths.TrashFolder, paths.MetadataFolder, trashRelPath+".json")
	metaOriginal := filepath.Join(userDir, paths.MetadataFolder, trashRelPath+".json")
	_ = moveFile(metaTrash, metaOriginal)
	return nil
}

// ListFiles returns all files in the Trash folder (excluding Thumbnails/Metadata).
// Paths are relative to userDir with forward slashes.
func ListFiles(userDir string) ([]string, error) {
	var files []string
	trashDir := filepath.Join(userDir, paths.TrashFolder)
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
		rel = paths.Normalize(rel)
		if strings.HasPrefix(rel, paths.TrashFolder+"/"+paths.ThumbnailsFolder+"/") ||
			strings.HasPrefix(rel, paths.TrashFolder+"/"+paths.MetadataFolder+"/") {
			return nil
		}
		files = append(files, rel)
		return nil
	})
	return files, err
}

// moveFile moves a file from src to dst, creating destination directories as needed.
// Returns nil if source doesn't exist.
func moveFile(src, dst string) error {
	if _, err := os.Stat(src); os.IsNotExist(err) {
		return nil
	}
	if err := os.MkdirAll(filepath.Dir(dst), 0755); err != nil {
		return err
	}
	return os.Rename(src, dst)
}

// TrashPrefix returns the trash folder prefix with trailing slash ("Trash/").
func TrashPrefix() string {
	return trashPrefix
}
