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
	"strings"

	"github.com/takecontrolsoft/sync_server/server/config"
	"github.com/takecontrolsoft/sync_server/server/store"
	"github.com/takecontrolsoft/sync_server/server/utils"
)

type loginRequest struct {
	User     string `json:"User"`
	Password string `json:"Password"`
}

type loginResponse struct {
	Token  string `json:"Token"`
	UserId string `json:"UserId"`
}

type registerRequest struct {
	User     string `json:"User"`
	Password string `json:"Password"`
}

// LoginHandler validates user/password and returns a session token.
// POST body: { "User": "", "Password": "" } -> { "Token": "" } or 401.
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	var req loginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RenderError(w, err, http.StatusBadRequest)
		return
	}
	if req.User == "" || req.Password == "" {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	if !store.VerifyUser(req.User, req.Password) {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	userId := store.GetUserIdByUsername(req.User)
	if userId == "" {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	token, err := store.CreateToken(userId)
	if err != nil {
		utils.RenderError(w, err, http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(loginResponse{Token: token, UserId: userId})
}

// RegisterHandler creates a user. POST body: { "User": "", "Password": "" }.
func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	var req registerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RenderError(w, err, http.StatusBadRequest)
		return
	}
	if req.User == "" || req.Password == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	userId, err := store.CreateUser(req.User, req.Password)
	if err != nil {
		utils.RenderError(w, err, http.StatusInternalServerError)
		return
	}
	token, _ := store.CreateToken(userId)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(loginResponse{Token: token, UserId: userId})
}

// ResolveToUserId returns the folder name used for storage path (UploadDirectory/<this>/deviceId).
// Always returns normalized (lowercase) email so the main folder is the user's email.
// When auth DB is set: User in requests can be username (email) or userId; we resolve to email for path.
// When auth DB is not set: returns lowercase user as-is.
func ResolveToUserId(user string) string {
	if user == "" {
		return ""
	}
	if config.AuthDBPath != "" && store.UserIdExists(user) {
		email := store.GetUsernameByUserId(user)
		if email != "" {
			return strings.ToLower(email)
		}
	}
	return strings.ToLower(user)
}

