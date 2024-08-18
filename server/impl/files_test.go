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
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-errors/errors"
	"github.com/takecontrolsoft/go_multi_log/logger"
	"github.com/takecontrolsoft/sync_server/server/config"
	"github.com/takecontrolsoft/sync_server/server/utils"
)

func init() {
	config.InitFromEnvVariables()
}

func TestGetFiles(t *testing.T) {
	userName := "Desi"
	deviceId := "AFA33F68-3E48-5459-B564-20D03E3F6035"
	folder := "2024/7"
	body, err := postFileForm(userName, deviceId, folder)
	if err != nil {
		t.Fatal(err)
	}

	for i := range body {
		f := body[i]
		fmt.Println(f)
	}
}

func postFileForm(userName string, deviceId string, folder string) ([]string, error) {
	body := folderData{UserData: userData{User: userName, DeviceId: deviceId}, Folder: folder}
	r, err := utils.JsonReaderFactory(body)
	if err != nil {
		logger.Error(err)
		return nil, err
	}
	req, err := http.NewRequest("POST", "/files", r)
	if err != nil {
		logger.Error(err)
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(GetFilesHandler)
	handler.ServeHTTP(rr, req)

	if rr.Code == http.StatusOK {
		result := []string{}
		if err := json.NewDecoder(rr.Body).Decode(&result); err != nil {
			logger.Error(err)
			return nil, err
		}
		return result, nil
	} else {
		return nil, &RequestError{
			StatusCode: rr.Code,
			Err:        errors.New(rr.Body),
		}
	}
}
