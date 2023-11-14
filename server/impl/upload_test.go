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
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
)

func TestHealthCheckHandler(t *testing.T) {
	r, w := io.Pipe()
	m := multipart.NewWriter(w)
	name := "/video.mp4"
	go writeAsync(w, m, name)
	req, err := http.NewRequest("POST", "/upload", r)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", m.FormDataContentType())

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(UploadHandler)

	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
}

func writeAsync(w *io.PipeWriter, m *multipart.Writer, fn string) {
	defer w.Close()
	defer m.Close()
	f := filepath.Base(fn)
	part, err := m.CreateFormFile(f, f)
	if err != nil {
		return
	}
	file, err := os.Open(fn)
	if err != nil {
		return
	}
	defer file.Close()
	if _, err = io.Copy(part, file); err != nil {
		return
	}
}
