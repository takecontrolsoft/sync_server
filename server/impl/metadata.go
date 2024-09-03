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

	et, err := exiftool.NewExiftool()
	if err != nil {
		return "", err
	}
	defer et.Close()

	fileInfos := et.ExtractMetadata(filePath)
	outputJson, err := json.Marshal(fileInfos)
	if err != nil {
		return "", err
	}
	err = os.WriteFile(metadataPath, outputJson, 0644)
	if err != nil {
		return "", err
	}
	return metadataPath, nil
}
