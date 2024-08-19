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
	"image"
	"image/png"
	"net/http"
	"path/filepath"

	"github.com/go-errors/errors"
	"github.com/takecontrolsoft/sync_server/server/config"
	"github.com/takecontrolsoft/sync_server/server/utils"
)

type fileData struct {
	UserData userData
	File     string
}

func GetImageHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		var files = make([]string, 0)
		var result fileData
		if err := json.NewDecoder(r.Body).Decode(&result); err != nil {
			utils.RenderError(w, errors.Errorf("$Required json input {UserData: { User: '', DeviceId: ''}, 	File: ''}"), http.StatusBadRequest)
			return
		}
		userName := result.UserData.User
		deviceId := result.UserData.DeviceId
		file := result.File
		userDirName := filepath.Join(config.UploadDirectory, userName, deviceId)
		filePath := filepath.Join(userDirName, file)
		src, err := utils.GetImageFromFilePath(filePath)
		if err != nil {
			utils.RenderError(w, err, http.StatusBadRequest)
			return
		}
		// Set the expected size that you want:
		dst := image.NewRGBA(image.Rect(0, 0, src.Bounds().Max.X/2, src.Bounds().Max.Y/2))

		// Encode to `output`:
		png.Encode(w, dst)

		w.Header().Set("Content-Type", "image/png")
		w.WriteHeader(http.StatusOK)

		if err := json.NewEncoder(w).Encode(files); err != nil {
			if utils.RenderIfError(err, w, http.StatusInternalServerError) {
				return
			}
		}
	}
}
