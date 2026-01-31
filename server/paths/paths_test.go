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

package paths

import (
	"path/filepath"
	"testing"
)

func TestNormalize(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"foo\\bar\\baz", "foo/bar/baz"},
		{"foo/bar/baz", "foo/bar/baz"},
		{"", ""},
		{"file.txt", "file.txt"},
	}
	for _, tt := range tests {
		got := Normalize(tt.input)
		if got != tt.expected {
			t.Errorf("Normalize(%q) = %q, want %q", tt.input, got, tt.expected)
		}
	}
}

func TestThumbnailBasePath(t *testing.T) {
	userDir := filepath.Join("uploads", "user1", "device1")

	tests := []struct {
		file     string
		expected string
	}{
		{"2024/01/photo.jpg", filepath.Join(userDir, ThumbnailsFolder, "2024/01/photo.jpg")},
		{"Trash/2024/01/photo.jpg", filepath.Join(userDir, TrashFolder, ThumbnailsFolder, "2024/01/photo.jpg")},
		{"trash/2024/01/photo.jpg", filepath.Join(userDir, TrashFolder, ThumbnailsFolder, "2024/01/photo.jpg")},
	}
	for _, tt := range tests {
		got := ThumbnailBasePath(userDir, tt.file)
		if got != tt.expected {
			t.Errorf("ThumbnailBasePath(%q, %q) = %q, want %q", userDir, tt.file, got, tt.expected)
		}
	}
}

func TestMetadataPath(t *testing.T) {
	userDir := filepath.Join("uploads", "user1", "device1")

	tests := []struct {
		file     string
		expected string
	}{
		{"2024/01/photo.jpg", filepath.Join(userDir, MetadataFolder, "2024/01/photo.jpg.json")},
		{"Trash/2024/01/photo.jpg", filepath.Join(userDir, TrashFolder, MetadataFolder, "2024/01/photo.jpg.json")},
	}
	for _, tt := range tests {
		got := MetadataPath(userDir, tt.file)
		if got != tt.expected {
			t.Errorf("MetadataPath(%q, %q) = %q, want %q", userDir, tt.file, got, tt.expected)
		}
	}
}

func TestShouldSkipInFolderListing(t *testing.T) {
	tests := []struct {
		rel      string
		expected bool
	}{
		{"Trash", true},
		{"Trash/2024", true},
		{"Thumbnails", true},
		{"Thumbnails/2024", true},
		{"Metadata", true},
		{"Metadata/2024", true},
		{"2024/01", false},
		{"2024/Thumbnails/photo.jpg", true},
		{"photo.jpg", false},
	}
	for _, tt := range tests {
		got := ShouldSkipInFolderListing(tt.rel)
		if got != tt.expected {
			t.Errorf("ShouldSkipInFolderListing(%q) = %v, want %v", tt.rel, got, tt.expected)
		}
	}
}

func TestIsImagePath(t *testing.T) {
	tests := []struct {
		path     string
		expected bool
	}{
		{"photo.jpg", true},
		{"photo.JPEG", true},
		{"photo.png", true},
		{"photo.gif", true},
		{"photo.webp", true},
		{"photo.heic", true},
		{"video.mp4", false},
		{"document.pdf", false},
		{"noextension", false},
		{"", false},
	}
	for _, tt := range tests {
		got := IsImagePath(tt.path)
		if got != tt.expected {
			t.Errorf("IsImagePath(%q) = %v, want %v", tt.path, got, tt.expected)
		}
	}
}

func TestIsUnderMetadata(t *testing.T) {
	tests := []struct {
		rel      string
		expected bool
	}{
		{"Metadata", true},
		{"Metadata/2024", true},
		{"2024/Metadata/file.json", true},
		{"2024/01/photo.jpg", false},
	}
	for _, tt := range tests {
		got := IsUnderMetadata(tt.rel)
		if got != tt.expected {
			t.Errorf("IsUnderMetadata(%q) = %v, want %v", tt.rel, got, tt.expected)
		}
	}
}
