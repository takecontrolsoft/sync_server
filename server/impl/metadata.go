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
	"os"
	"path/filepath"

	"github.com/barasher/go-exiftool"
	"github.com/takecontrolsoft/sync_server/server/config"
)

func ExtractMetadata(userName string, deviceId string, file string) (string, error) {

	userDirName := filepath.Join(config.UploadDirectory, userName, deviceId)
	metadataPath := filepath.Join(userDirName, "Metadata", fmt.Sprintf("%s.json", file))
	filePath := filepath.Join(userDirName, file)

	var et *exiftool.Exiftool
	var err error
	if config.BinDirectory != "" {
		et, err = exiftool.NewExiftool(exiftool.SetExiftoolBinaryPath(config.ExiftoolBinary()))
	} else {
		et, err = exiftool.NewExiftool()
	}
	if err != nil {
		return "", err
	}
	defer et.Close()

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
