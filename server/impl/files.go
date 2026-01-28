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
	"net/http"
	"os"
	"path/filepath"

	"github.com/go-errors/errors"
	"github.com/takecontrolsoft/sync_server/server/config"
	"github.com/takecontrolsoft/sync_server/server/utils"
)

type folderData struct {
	UserData userData
	Folder   string
}

func GetFilesHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		var files = make([]string, 0)
		var result folderData
		if err := json.NewDecoder(r.Body).Decode(&result); err != nil {
			utils.RenderError(w, errors.Errorf("$Required json input {UserData: { User: '', DeviceId: ''}, Folder: ''}"), http.StatusBadRequest)
			return
		}
		userFromClient := result.UserData.User
		deviceId := result.UserData.DeviceId
		folder := result.Folder
		userId := ResolveToUserId(userFromClient)
		if userId == "" {
			userId = userFromClient
		}
		userDirName := filepath.Join(config.UploadDirectory, userId, deviceId)
		if folder == TrashFolder {
			files, _ = ListTrashFiles(userDirName)
		} else {
			dirName := filepath.Join(userDirName, folder)
			entries, err := os.ReadDir(dirName)
			if err == nil {
				for _, entry := range entries {
					if !entry.IsDir() {
						file := filepath.Join(folder, entry.Name())
						files = append(files, file)
					}
				}
			}
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		if err := json.NewEncoder(w).Encode(files); err != nil {
			if utils.RenderIfError(err, w, http.StatusInternalServerError) {
				return
			}
		}
	}
}

