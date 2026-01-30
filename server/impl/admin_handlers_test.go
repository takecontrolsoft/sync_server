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
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/takecontrolsoft/sync_server/server/config"
	"github.com/takecontrolsoft/sync_server/server/utils"
)

func init() {
	// Allow tests to run without env set (e.g. go test ./server/impl/ -run "...")
	if os.Getenv(config.UploadPathVariable) == "" {
		_ = os.Setenv(config.UploadPathVariable, os.TempDir())
	}
	if os.Getenv(config.PortVariable) == "" {
		_ = os.Setenv(config.PortVariable, "8080")
	}
	config.InitFromEnvVariables()
}

func TestRegenerateThumbnailsHandler_methodNotAllowed(t *testing.T) {
	body := struct {
		UserData userData `json:"UserData"`
	}{UserData: userData{User: "test@example.com", DeviceId: "device1"}}
	r, _ := utils.JsonReaderFactory(body)
	req := httptest.NewRequest(http.MethodGet, "/regenerate-thumbnails", r)
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	RegenerateThumbnailsHandler(rr, req)
	if rr.Code != http.StatusMethodNotAllowed {
		t.Errorf("RegenerateThumbnailsHandler GET: got status %d, want %d", rr.Code, http.StatusMethodNotAllowed)
	}
}

func TestRegenerateThumbnailsHandler_badRequest(t *testing.T) {
	tests := []struct {
		name string
		body string
	}{
		{"empty body", ""},
		{"missing User", `{"UserData":{"User":"","DeviceId":"dev1"}}`},
		{"missing DeviceId", `{"UserData":{"User":"u","DeviceId":""}}`},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/regenerate-thumbnails", bytes.NewBufferString(tt.body))
			req.Header.Set("Content-Type", "application/json")
			rr := httptest.NewRecorder()
			RegenerateThumbnailsHandler(rr, req)
			if rr.Code != http.StatusBadRequest {
				t.Errorf("got status %d, want %d", rr.Code, http.StatusBadRequest)
			}
		})
	}
}

func TestRegenerateThumbnailsHandler_ok(t *testing.T) {
	tmp := t.TempDir()
	restore := config.UploadDirectory
	config.UploadDirectory = tmp
	defer func() { config.UploadDirectory = restore }()

	userDir := filepath.Join(tmp, "testuser", "device1")
	if err := os.MkdirAll(userDir, 0755); err != nil {
		t.Fatal(err)
	}

	body := struct {
		UserData userData `json:"UserData"`
	}{UserData: userData{User: "testuser", DeviceId: "device1"}}
	r, _ := utils.JsonReaderFactory(body)
	req := httptest.NewRequest(http.MethodPost, "/regenerate-thumbnails", r)
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	RegenerateThumbnailsHandler(rr, req)
	if rr.Code != http.StatusOK {
		t.Errorf("got status %d, want %d", rr.Code, http.StatusOK)
	}
	var res map[string]int
	if err := json.NewDecoder(rr.Body).Decode(&res); err != nil {
		t.Fatal(err)
	}
	if _, ok := res["Regenerated"]; !ok {
		t.Errorf("response missing Regenerated: %v", res)
	}
}

func TestCleanOrphanThumbnailsHandler_methodNotAllowed(t *testing.T) {
	body := struct {
		UserData userData `json:"UserData"`
	}{UserData: userData{User: "u", DeviceId: "d"}}
	r, _ := utils.JsonReaderFactory(body)
	req := httptest.NewRequest(http.MethodGet, "/clean-orphan-thumbnails", r)
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	CleanOrphanThumbnailsHandler(rr, req)
	if rr.Code != http.StatusMethodNotAllowed {
		t.Errorf("got status %d, want %d", rr.Code, http.StatusMethodNotAllowed)
	}
}

func TestCleanOrphanThumbnailsHandler_badRequest(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/clean-orphan-thumbnails", bytes.NewBufferString(`{}`))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	CleanOrphanThumbnailsHandler(rr, req)
	if rr.Code != http.StatusBadRequest {
		t.Errorf("got status %d, want %d", rr.Code, http.StatusBadRequest)
	}
}

func TestCleanOrphanThumbnailsHandler_ok(t *testing.T) {
	tmp := t.TempDir()
	restore := config.UploadDirectory
	config.UploadDirectory = tmp
	defer func() { config.UploadDirectory = restore }()

	userDir := filepath.Join(tmp, "user1", "dev1")
	if err := os.MkdirAll(userDir, 0755); err != nil {
		t.Fatal(err)
	}

	body := struct {
		UserData userData `json:"UserData"`
	}{UserData: userData{User: "user1", DeviceId: "dev1"}}
	r, _ := utils.JsonReaderFactory(body)
	req := httptest.NewRequest(http.MethodPost, "/clean-orphan-thumbnails", r)
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	CleanOrphanThumbnailsHandler(rr, req)
	if rr.Code != http.StatusOK {
		t.Errorf("got status %d, want %d", rr.Code, http.StatusOK)
	}
	var res map[string]int
	if err := json.NewDecoder(rr.Body).Decode(&res); err != nil {
		t.Fatal(err)
	}
	if _, ok := res["Removed"]; !ok {
		t.Errorf("response missing Removed: %v", res)
	}
}

func TestRunDocumentDetectionHandler_methodNotAllowed(t *testing.T) {
	body := struct {
		UserData userData `json:"UserData"`
	}{UserData: userData{User: "u", DeviceId: "d"}}
	r, _ := utils.JsonReaderFactory(body)
	req := httptest.NewRequest(http.MethodGet, "/run-document-detection", r)
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	RunDocumentDetectionHandler(rr, req)
	if rr.Code != http.StatusMethodNotAllowed {
		t.Errorf("got status %d, want %d", rr.Code, http.StatusMethodNotAllowed)
	}
}

func TestRunDocumentDetectionHandler_badRequest(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/run-document-detection", bytes.NewBufferString(`{"UserData":{"User":"","DeviceId":"d"}}`))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	RunDocumentDetectionHandler(rr, req)
	if rr.Code != http.StatusBadRequest {
		t.Errorf("got status %d, want %d", rr.Code, http.StatusBadRequest)
	}
}

func TestRunDocumentDetectionHandler_ok_emptyDir(t *testing.T) {
	tmp := t.TempDir()
	restore := config.UploadDirectory
	config.UploadDirectory = tmp
	defer func() { config.UploadDirectory = restore }()

	userDir := filepath.Join(tmp, "nobody", "dev1")
	if err := os.MkdirAll(userDir, 0755); err != nil {
		t.Fatal(err)
	}

	body := struct {
		UserData userData `json:"UserData"`
	}{UserData: userData{User: "nobody", DeviceId: "dev1"}}
	r, _ := utils.JsonReaderFactory(body)
	req := httptest.NewRequest(http.MethodPost, "/run-document-detection", r)
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	RunDocumentDetectionHandler(rr, req)
	if rr.Code != http.StatusOK {
		t.Errorf("got status %d, want %d", rr.Code, http.StatusOK)
	}
	var res map[string]int
	if err := json.NewDecoder(rr.Body).Decode(&res); err != nil {
		t.Fatal(err)
	}
	if res["Moved"] != 0 {
		t.Errorf("Moved: got %d, want 0", res["Moved"])
	}
}

func TestRunDocumentDetectionHandler_ok_userDirNotExist(t *testing.T) {
	tmp := t.TempDir()
	restore := config.UploadDirectory
	config.UploadDirectory = tmp
	defer func() { config.UploadDirectory = restore }()

	// Do not create user dir - handler should return 200 with Moved: 0
	body := struct {
		UserData userData `json:"UserData"`
	}{UserData: userData{User: "nonexistent", DeviceId: "dev1"}}
	r, _ := utils.JsonReaderFactory(body)
	req := httptest.NewRequest(http.MethodPost, "/run-document-detection", r)
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	RunDocumentDetectionHandler(rr, req)
	if rr.Code != http.StatusOK {
		t.Errorf("got status %d, want %d", rr.Code, http.StatusOK)
	}
	var res map[string]int
	if err := json.NewDecoder(rr.Body).Decode(&res); err != nil {
		t.Fatal(err)
	}
	if res["Moved"] != 0 {
		t.Errorf("Moved when dir missing: got %d, want 0", res["Moved"])
	}
}
