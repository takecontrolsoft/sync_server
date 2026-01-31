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

// Package paths provides path constants and utilities for the sync server storage structure.
package paths

import (
	"os"
	"path/filepath"
	"strings"
)

// Storage folder constants
const (
	TrashFolder      = "Trash"
	ThumbnailsFolder = "Thumbnails"
	MetadataFolder   = "Metadata"
)

// Normalize converts backslashes to forward slashes for cross-platform consistency.
// All API responses and internal path comparisons use forward slashes.
func Normalize(path string) string {
	return filepath.ToSlash(path)
}

// TrimLeadingSeparator removes a leading OS path separator if present.
func TrimLeadingSeparator(path string) string {
	return strings.TrimPrefix(path, string(os.PathSeparator))
}

// ThumbnailBasePath returns the thumbnail storage path for a given file.
// Handles Trash paths: Trash/2024/file.jpg -> userDir/Trash/Thumbnails/2024/file.jpg
func ThumbnailBasePath(userDir, file string) string {
	file = Normalize(file)
	file = strings.TrimSpace(file)
	parts := strings.Split(file, "/")
	if len(parts) > 0 && strings.EqualFold(parts[0], TrashFolder) {
		rest := strings.Join(parts[1:], "/")
		return filepath.Join(userDir, TrashFolder, ThumbnailsFolder, rest)
	}
	return filepath.Join(userDir, ThumbnailsFolder, file)
}

// MetadataPath returns the metadata JSON storage path for a given file.
// Handles Trash paths: Trash/2024/file.jpg -> userDir/Trash/Metadata/2024/file.jpg.json
func MetadataPath(userDir, file string) string {
	file = Normalize(file)
	file = strings.TrimSpace(file)
	parts := strings.Split(file, "/")
	if len(parts) > 0 && strings.EqualFold(parts[0], TrashFolder) {
		rest := strings.Join(parts[1:], "/")
		return filepath.Join(userDir, TrashFolder, MetadataFolder, rest+".json")
	}
	return filepath.Join(userDir, MetadataFolder, file+".json")
}

// IsUnderMetadata returns true if the path is inside a Metadata folder.
func IsUnderMetadata(rel string) bool {
	norm := Normalize(TrimLeadingSeparator(rel))
	return norm == MetadataFolder ||
		strings.HasPrefix(norm, MetadataFolder+"/") ||
		strings.Contains(norm, "/"+MetadataFolder+"/")
}

// ShouldSkipInFolderListing returns true if the path should be excluded from folder/file listings.
// Excludes: Trash, Thumbnails, Metadata directories.
func ShouldSkipInFolderListing(rel string) bool {
	norm := Normalize(TrimLeadingSeparator(rel))
	if norm == TrashFolder || strings.HasPrefix(norm, TrashFolder+"/") {
		return true
	}
	if norm == ThumbnailsFolder ||
		strings.HasPrefix(norm, ThumbnailsFolder+"/") ||
		strings.Contains(norm, "/"+ThumbnailsFolder+"/") {
		return true
	}
	return IsUnderMetadata(rel)
}

// IsImagePath returns true if the file extension is a common image type (case-insensitive).
func IsImagePath(path string) bool {
	ext := strings.ToLower(filepath.Ext(path))
	if ext == "" {
		return false
	}
	ext = ext[1:] // drop leading dot
	switch ext {
	case "jpg", "jpeg", "png", "gif", "bmp", "webp", "heic":
		return true
	}
	return false
}
