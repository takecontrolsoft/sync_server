/*
	Copyright 2023 Take Control - Software & Infrastructure

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
	"net/http"

	"github.com/go-errors/errors"
)

// An error for invalid file type, which is not allowed to be uploaded.
func InvalidFileTypeUploaded(fileType string) error {
	return errors.Errorf("File type '%s' is not allowed to be uploaded.", fileType).Err
}

// An error for empty storage path.
var FileSizeExceeded = errors.Errorf("Maximum file size for uploaded files exceeded.").Err

// An error for missing authorized user.
var MissingUser = errors.Errorf("The user is not authorized.").Err

type RequestError struct {
	StatusCode int

	Err error
}

func (r *RequestError) Error() string {
	return r.Err.Error()
}

func (r *RequestError) ServiceUnavailable() bool {
	return r.StatusCode == http.StatusServiceUnavailable
}

func (r *RequestError) InternalServerError() bool {
	return r.StatusCode == http.StatusInternalServerError
}

func (r *RequestError) BadRequest() bool {
	return r.StatusCode == http.StatusBadRequest
}
