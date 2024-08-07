/* Copyright 2023 Take Control - Software & Infrastructure

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
	"path/filepath"
	"strings"

	"github.com/takecontrolsoft/sync_server/server/config"
	"github.com/takecontrolsoft/sync_server/server/utils"
)

type data struct {
	User     string
	DeviceId string
}

func GetFoldersHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		var folders []string
		var result data
		if err := json.NewDecoder(r.Body).Decode(&result); err != nil {
			utils.RenderError(w, err, http.StatusBadRequest)
			return
		}
		userName := result.User
		deviceId := result.DeviceId
		dirName := filepath.Join(config.UploadDirectory, userName, deviceId)

		err := filepath.WalkDir(dirName, func(path string, d fs.DirEntry, err error) error {

			if d.IsDir() && deviceId != d.Name() {
				folders = append(
					folders,
					strings.Replace(path+"/", dirName, "", 1),
				)
			}
			return err
		})
		if utils.RenderIfError(err, w, http.StatusInternalServerError) {
			return
		}
		jData, err := json.Marshal(folders)
		if utils.RenderIfError(err, w, http.StatusInternalServerError) {
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		w.Write(jData)

	}
}
