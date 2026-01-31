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

// Package auth provides authentication utilities including user ID resolution.
package auth

import (
	"strings"

	"github.com/takecontrolsoft/sync_server/server/config"
	"github.com/takecontrolsoft/sync_server/server/store"
)

// ResolveUserId returns the folder name used for storage path (UploadDirectory/<this>/deviceId).
// Always returns normalized (lowercase) email so the main folder is the user's email.
// When auth DB is set: User in requests can be username (email) or userId; we resolve to email for path.
// When auth DB is not set: returns lowercase user as-is.
func ResolveUserId(user string) string {
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
