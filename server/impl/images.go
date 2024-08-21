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
	"image/png"
	"net/http"
	"os"
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
		var result fileData
		if err := json.NewDecoder(r.Body).Decode(&result); err != nil {
			utils.RenderError(w, errors.Errorf("$Required json input {UserData: { User: '', DeviceId: ''}, 	File: ''}"), http.StatusBadRequest)
			return
		}
		userName := result.UserData.User
		deviceId := result.UserData.DeviceId
		file := result.File
		userDirName := filepath.Join(config.UploadDirectory, userName, deviceId)
		var thumbnailPath = filepath.Join(userDirName, "Thumbnails", file)

		src, err := utils.GetImageFromFilePath(thumbnailPath)

		if err != nil {
			utils.RenderError(w, err, http.StatusInternalServerError)
			return
		}
		png.Encode(w, src)

		w.Header().Set("Content-Type", "image/png")
		w.WriteHeader(http.StatusOK)
	}
}

func BuildThumbnail(userName string, deviceId string, file string) (string, error) {
	userDirName := filepath.Join(config.UploadDirectory, userName, deviceId)
	thumbnailPath := filepath.Join(userDirName, "Thumbnails", file)
	filePath := filepath.Join(userDirName, file)
	src, err := utils.GetImageFromFilePath(filePath)
	if err != nil {
		return "", err
	}

	rgba_src := utils.ImageToRGBA(src)
	resized := utils.ResizeImage(rgba_src, 90)
	err = os.MkdirAll(filepath.Dir(thumbnailPath), os.ModePerm)
	if err != nil {
		return "", err
	}
	f, err := os.Create(thumbnailPath)
	if err != nil {
		return "", err
	}
	defer f.Close()
	png.Encode(f, resized)
	return thumbnailPath, nil
}
