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
	"io/fs"
	"os"
	"path/filepath"

	"github.com/barasher/go-exiftool"
	"github.com/takecontrolsoft/go_multi_log/logger"
	"github.com/takecontrolsoft/sync_server/server/config"
	"github.com/takecontrolsoft/sync_server/server/mediatypes"
	"github.com/takecontrolsoft/sync_server/server/utils"
)

// ExtractMetadata extracts the metadata for a single stored file and writes it
// as JSON under the "Metadata" folder, overwriting any existing metadata file.
func ExtractMetadata(userName string, deviceId string, file string) (string, error) {
	et, err := newExiftool()
	if err != nil {
		return "", err
	}
	defer et.Close()

	return extractMetadataWith(et, userName, deviceId, file)
}

// newExiftool opens an exiftool instance, using the bundled binary next to the
// executable when BinDirectory is configured.
func newExiftool() (*exiftool.Exiftool, error) {
	if config.BinDirectory != "" {
		return exiftool.NewExiftool(exiftool.SetExiftoolBinaryPath(config.ExiftoolBinary()))
	}
	return exiftool.NewExiftool()
}

// extractMetadataWith performs the extraction using an already opened exiftool
// instance so that it can be reused across many files (e.g. bulk regeneration).
func extractMetadataWith(et *exiftool.Exiftool, userName string, deviceId string, file string) (string, error) {
	userDirName := filepath.Join(config.UploadDirectory, userName, deviceId)
	// Use MetadataPath so Trash files get metadata under Trash/Metadata/, not Metadata/Trash/.
	metadataPath := MetadataPath(userDirName, file)
	filePath := filepath.Join(userDirName, file)

	fileInfos := et.ExtractMetadata(filePath)
	outputJson, err := json.Marshal(fileInfos)
	if err != nil {
		return "", err
	}
	dir := filepath.Dir(metadataPath)
	err = os.MkdirAll(dir, 0755)
	if err != nil {
		return "", err
	}
	err = os.WriteFile(metadataPath, outputJson, 0644)
	if err != nil {
		return "", err
	}
	return metadataPath, nil
}

// GetOrientationFromMetadata reads the EXIF Orientation (1-8) from a metadata JSON file.
// Returns 1 (normal) if the file is missing, invalid, or Orientation is absent.
func GetOrientationFromMetadata(metadataPath string) int {
	data, err := os.ReadFile(metadataPath)
	if err != nil {
		return 1
	}
	var fileInfos []struct {
		Fields map[string]interface{} `json:"Fields"`
	}
	if err := json.Unmarshal(data, &fileInfos); err != nil {
		return 1
	}
	if len(fileInfos) == 0 || fileInfos[0].Fields == nil {
		return 1
	}
	o := fileInfos[0].Fields["Orientation"]
	if o == nil {
		return 1
	}
	switch v := o.(type) {
	case float64:
		orient := int(v)
		if orient >= 1 && orient <= 8 {
			return orient
		}
	case int:
		if v >= 1 && v <= 8 {
			return v
		}
	}
	return 1
}

// RegenResult holds the counters produced by a regeneration run.
type RegenResult struct {
	MetadataOK      int
	MetadataFailed  int
	ThumbnailOK     int
	ThumbnailFailed int
	Skipped         int
}

// RegenerateForDevice walks every stored file for the given user and device and
// regenerates its Metadata JSON file and its Thumbnail, overwriting any existing
// ones. The generated "Metadata" and "Thumbnails" folders are skipped so they
// are not processed as input.
func RegenerateForDevice(userName string, deviceId string) (RegenResult, error) {
	var res RegenResult

	root := config.UploadDirectory
	if root == "" {
		return res, fmt.Errorf("upload directory is not configured")
	}
	deviceRoot := filepath.Join(root, userName, deviceId)
	if info, statErr := os.Stat(deviceRoot); statErr != nil || !info.IsDir() {
		return res, fmt.Errorf("device directory not found: %s", deviceRoot)
	}

	et, err := newExiftool()
	if err != nil {
		return res, err
	}
	defer et.Close()

	walkErr := filepath.WalkDir(deviceRoot, func(path string, d fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if d.IsDir() {
			// Do not process the generated folders as input.
			if d.Name() == "Metadata" || d.Name() == "Thumbnails" {
				return filepath.SkipDir
			}
			return nil
		}
		// file is the path relative to the device directory, e.g. <year>/<month>/<file>.
		file, relErr := filepath.Rel(deviceRoot, path)
		if relErr != nil {
			return nil
		}

		// Metadata.
		if _, mErr := extractMetadataWith(et, userName, deviceId, file); mErr != nil {
			logger.ErrorF("Regenerating metadata failed for %s: %v", file, mErr)
			res.MetadataFailed++
		} else {
			res.MetadataOK++
		}

		// Thumbnail, dispatched by detected media type.
		mediaType, mtErr := utils.GetMediaTypeFromFile(path)
		if mtErr != nil {
			logger.ErrorF("Detecting media type failed for %s: %v", file, mtErr)
		}
		var tErr error
		switch mediaType {
		case mediatypes.Video:
			_, tErr = BuildVideoThumbnail(userName, deviceId, file)
		case mediatypes.Image:
			_, tErr = BuildImageThumbnail(userName, deviceId, file)
		case mediatypes.Audio:
			_, tErr = BuildAudioThumbnail(userName, deviceId, file)
		default:
			res.Skipped++
			return nil
		}
		if tErr != nil {
			logger.ErrorF("Regenerating thumbnail failed for %s: %v", file, tErr)
			res.ThumbnailFailed++
		} else {
			res.ThumbnailOK++
		}
		return nil
	})

	logger.InfoF("Regeneration complete for %s/%s: metadata %d ok / %d failed, thumbnails %d ok / %d failed, %d skipped",
		userName, deviceId, res.MetadataOK, res.MetadataFailed, res.ThumbnailOK, res.ThumbnailFailed, res.Skipped)
	return res, walkErr
}
