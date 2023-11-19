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
	"path/filepath"
	"testing"
	"testing/fstest"
	"time"

	"github.com/go-errors/errors"
	"github.com/stretchr/testify/assert"
	"github.com/takecontrolsoft/sync/server/config"
	"github.com/takecontrolsoft/sync/server/utils"
)

func init() {
	config.InitFromEnvVariables()
}

func TestValidFiles(t *testing.T) {
	name := "C:\\Users\\desis\\Pictures\\veranda.png"
	field := utils.GenerateRandomString(5)
	err := postUpload(field, name, t)
	if err != nil {
		t.Fatal(err)
	}
}

func TestInvalidFileNames_Uploaded(t *testing.T) {
	name := "C:\\Users\\desis\\Pictures\\veranda.png"
	field := ")**??/\\\\//<>***_+"
	err := postUpload(field, name, t)
	if err != nil {
		t.Fatal(err)
	}
}

func TestInvalidFileTypes_Failed(t *testing.T) {
	name := "video.mp4"
	field := ")**??/\\\\//<>***_+"
	assert := assert.New(t)
	err := postUpload(field, name, t)
	if err != nil {
		re, ok := err.(*RequestError)
		if ok {
			assert.True(re.BadRequest(), "Expected bad request error.")
		} else {
			t.Fatal(err)
		}
	}
}

func postUpload(field string, name string, t *testing.T) error {
	r, w := io.Pipe()
	m := multipart.NewWriter(w)
	go writeAsync(w, m, field, name)
	req, err := http.NewRequest("POST", "/upload", r)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", m.FormDataContentType())

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(UploadHandler)

	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		return &RequestError{
			StatusCode: status,
			Err:        errors.New(rr.Body),
		}
	}
	return nil
}

func writeAsync(w *io.PipeWriter, m *multipart.Writer, field string, fn string) {
	defer w.Close()
	defer m.Close()
	fs := fstest.MapFS{
		fn: {
			Data:    getFileBytes(fn),
			Mode:    0,
			ModTime: time.Time{},
			Sys:     nil,
		},
	}
	f := filepath.Base(fn)
	part, err := m.CreateFormFile(field, f)
	if err != nil {
		return
	}
	file, err := fs.Open(fn)
	if err != nil {
		return
	}
	defer file.Close()
	if _, err = io.Copy(part, file); err != nil {
		return
	}
}

func getFileBytes(filename string) []byte {
	ext := filepath.Ext(filename)
	result := ""
	switch ext {
	case "text/plain; charset=utf-8", "audio/aac", "image/jpeg", "image/jpg", "image/gif", "image/png", "application/pdf":
		result = ""
		break
	case ".mp4":
		result = "\x00\x00\x00\x1cftypXAVC\x01\x00\x1f\xffXAVCmp42iso2\x00\x00\x00\x94uuidPROF!\xd2Oλ\x88i\\\xfa\xc9\xc7@\x00\x00\x00\x00\x00\x00\x00\x03\x00\x00\x00\x14FPRF\x00\x00\x00\x00 \x00\x00\x00\x00\x00\x00\x00\x00\x00\x00,APRF\x00\x00\x00\x00\x00\x00\x00\x02twos\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x06\x00\x00\x00\x06\x00\x00\x00\xbb\x80\x00\x00\x00\x02\x00\x00\x004VPRF\x00\x00\x00\x00\x00\x00\x00\x01avc1\x01d\x003\x00\x03\x00\x02\x00\x01\x86\xa0\x00\x01\x86\xa0\x00\x19\x00\x00\x00\x19\x00\x00\x0f\x00\bp\x00\x01\x00\x01\x00\x00\x00\x01mdat\x00\x00\x00\x00_\x7f\xffP\x00\x1c\x01\x00\x00\xfd\x02\xe7\xf0\x01\x00\x10 \x01*\x1f\x00\x11\x00\x19\x00\x00\x00\x00\x00\x00\x00\x00\x06\x0e+4\x02S\x01\x01\f\x02\x01\x01\x01\x01\x00\x00\x83\x00\x00\f\x80\x00\x00\x02\xb5\xb3\x80\x01\x00\x02\"\x8f\x06\x0e+4\x02S\x01\x01\f\x02\x01\x01\x02\x01\x00\x00\x83\x00\x00~\x81\x00\x00\x10\x06\x0e+4\x04\x01\x01\v\x05\x10\x01\x01\x01\x02\x00\x00\x81\x15\x00\x02\x03 \x81\x01\x00\x01\x02\x81\t\x00\b\x00\x00\x00\x01\x00\x00\x00\x19\x81\n\x00\x02\a\b\x81\v\x00\x02\x03 \x81\f\x00\x02\x00d\x81\r\x00\x01\x012\x10\x00\x10\x06\x0e+4\x04\x01\x01\r\x04\x01\x01\x01\x01\b\x00\x002\x19\x00\x10\x06\x0e+4\x04\x01\x01\x06\x04\x01\x01\x01\x03\x03\x00\x002\x1a\x00\x10\x06\x0e+4\x04\x01\x01\x01\x04\x01\x01\x01\x02\x02\x00\x00\x06\x0e+4\x02S\x01\x01\f\x02\x01\x01\x7f\x01\x00\x00\x83\x00\x007\xe0\x00\x00\x10\x96i\b\x00Fx\x03\x1c Q\x00\x00\xf0\xc0\x11\x81\xe3\x00\x00\x01\x00\xe3\x01\x00\x04\x00\x00\x03 \xe3\x02\x00\x01\x01\xe3\x03\x00\x01\xff\xe3\x04\x00\bF \x19\x06\x01\x13H5\x00\x00\x02\xe7kkad\x00\x00\x00\x00\x00\x00\x05\xd9\x00 \xf1\x191202010401181013\xf0\x02\x00"
		break
	case ".exe":
		result = "MZP\x00\x02\x00\x00\x00\x04\x00\x0f\x00\xff\xff\x00\x00\xb8\x00\x00\x00\x00\x00\x00\x00@\x00\x1a\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x01\x00\x00"
		break
	default:
		result = ""
		break
	}
	return []byte(result)
}
