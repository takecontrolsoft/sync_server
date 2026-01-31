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

package impl

import (
	"encoding/json"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/takecontrolsoft/sync_server/server/auth"
	"github.com/takecontrolsoft/sync_server/server/config"
	"github.com/takecontrolsoft/sync_server/server/paths"
	"github.com/takecontrolsoft/sync_server/server/trash"
	"github.com/takecontrolsoft/sync_server/server/utils"
)

type moveToTrashData struct {
	UserData userData
	Files    []string
}

type restoreData struct {
	UserData userData
	Files    []string
}

func MoveToTrashHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	var result moveToTrashData
	if err := json.NewDecoder(r.Body).Decode(&result); err != nil {
		utils.RenderError(w, err, http.StatusBadRequest)
		return
	}
	userFromClient := result.UserData.User
	deviceId := result.UserData.DeviceId
	if userFromClient == "" || deviceId == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	userId := auth.ResolveUserId(userFromClient)
	if userId == "" {
		userId = userFromClient
	}
	userDir := filepath.Join(config.UploadDirectory, userId, deviceId)

	for _, file := range result.Files {
		if file == "" || strings.Contains(file, "..") {
			continue
		}
		if strings.HasPrefix(file, trash.TrashPrefix()) || file == paths.TrashFolder {
			continue
		}
		_ = trash.MoveWithAssociatedFiles(userDir, file)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

func RestoreHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	var result restoreData
	if err := json.NewDecoder(r.Body).Decode(&result); err != nil {
		utils.RenderError(w, err, http.StatusBadRequest)
		return
	}
	userFromClient := result.UserData.User
	deviceId := result.UserData.DeviceId
	if userFromClient == "" || deviceId == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	userId := auth.ResolveUserId(userFromClient)
	if userId == "" {
		userId = userFromClient
	}
	userDir := filepath.Join(config.UploadDirectory, userId, deviceId)

	for _, file := range result.Files {
		if file == "" || strings.Contains(file, "..") {
			continue
		}
		if !strings.HasPrefix(file, trash.TrashPrefix()) {
			continue
		}
		restorePath := strings.TrimPrefix(file, trash.TrashPrefix())
		_ = trash.RestoreFromTrash(userDir, restorePath)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}
