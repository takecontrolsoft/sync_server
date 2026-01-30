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

// Package store provides a small local SQLite database for user names and
// password hashes, used for login and registration.
package store

import (
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/takecontrolsoft/go_multi_log/logger"
	"golang.org/x/crypto/bcrypt"

	_ "modernc.org/sqlite"
)

const tokenExpiryHours = 24 * 7 // 7 days

var (
	db   *sql.DB
	once sync.Once
)

// Open opens the auth SQLite database at path (creates file and dirs if needed). Safe to call once.
func Open(path string) error {
	var err error
	once.Do(func() {
		if path == "" {
			return
		}
		dir := filepath.Dir(path)
		if err = os.MkdirAll(dir, 0755); err != nil {
			return
		}
		db, err = sql.Open("sqlite", path)
		if err != nil {
			return
		}
		err = initSchema()
	})
	return err
}

func initSchema() error {
	if _, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			id TEXT PRIMARY KEY,
			username TEXT UNIQUE NOT NULL,
			password_hash TEXT NOT NULL,
			created_at INTEGER NOT NULL
		);
	`); err != nil {
		return err
	}
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS sessions (
			token TEXT PRIMARY KEY,
			user_id TEXT NOT NULL,
			expires_at INTEGER NOT NULL
		);
	`)
	return err
}

// CreateUser adds a user with the given password (hashed with bcrypt). Returns userId (UUID). If user already exists, returns their id and nil.
func CreateUser(username, password string) (userId string, err error) {
	if db == nil || username == "" || password == "" {
		return "", nil
	}
	if id := GetUserIdByUsername(username); id != "" {
		return id, nil
	}
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	userId = hex.EncodeToString(b)
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	_, err = db.Exec(
		`INSERT INTO users (id, username, password_hash, created_at) VALUES (?, ?, ?, strftime('%s','now'))`,
		userId, username, string(hash),
	)
	if err != nil {
		return "", err
	}
	return userId, nil
}

// VerifyUser returns true if the username exists and the password matches.
func VerifyUser(username, password string) bool {
	if db == nil || username == "" || password == "" {
		return false
	}
	var hash string
	err := db.QueryRow(`SELECT password_hash FROM users WHERE username = ?`, username).Scan(&hash)
	if err == sql.ErrNoRows || err != nil {
		return false
	}
	err = bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// HasAnyUser returns true if at least one user exists (for bootstrap).
func HasAnyUser() (bool, error) {
	if db == nil {
		return false, nil
	}
	var n int
	err := db.QueryRow(`SELECT COUNT(*) FROM users`).Scan(&n)
	return n > 0, err
}

// BootstrapFromEnv creates the first user from SYNC_ADMIN_USER and SYNC_ADMIN_PASSWORD if the DB has no users.
func BootstrapFromEnv(adminUser, adminPassword string) error {
	if db == nil || adminUser == "" || adminPassword == "" {
		return nil
	}
	ok, err := HasAnyUser()
	if err != nil || ok {
		return err
	}
	logger.InfoF("Auth DB: bootstrapping first user from env: %s", adminUser)
	_, err = CreateUser(adminUser, adminPassword)
	return err
}

// GetUserIdByUsername returns the user's id (UUID) or empty string if not found.
func GetUserIdByUsername(username string) string {
	if db == nil || username == "" {
		return ""
	}
	var id string
	err := db.QueryRow(`SELECT id FROM users WHERE username = ?`, username).Scan(&id)
	if err == sql.ErrNoRows || err != nil {
		return ""
	}
	return id
}

// GetUsernameByUserId returns the username (email) for the given userId, or empty string if not found.
func GetUsernameByUserId(userId string) string {
	if db == nil || userId == "" {
		return ""
	}
	var username string
	err := db.QueryRow(`SELECT username FROM users WHERE id = ?`, userId).Scan(&username)
	if err == sql.ErrNoRows || err != nil {
		return ""
	}
	return username
}

// UserIdExists returns true if the id exists in users.
func UserIdExists(id string) bool {
	if db == nil || id == "" {
		return false
	}
	var n int
	err := db.QueryRow(`SELECT 1 FROM users WHERE id = ?`, id).Scan(&n)
	return err == nil && n == 1
}

// CreateToken creates a session token for the user (by userId) and returns it. Caller must store it.
func CreateToken(userId string) (string, error) {
	if db == nil || userId == "" {
		return "", nil
	}
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	token := hex.EncodeToString(b)
	expires := time.Now().Add(tokenExpiryHours * time.Hour).Unix()
	_, err := db.Exec(`INSERT INTO sessions (token, user_id, expires_at) VALUES (?, ?, ?)`, token, userId, expires)
	return token, err
}

// ValidateToken returns the user id if the token is valid and not expired, else empty string.
func ValidateToken(token string) string {
	if db == nil || token == "" {
		return ""
	}
	var userId string
	var expiresAt int64
	err := db.QueryRow(`SELECT user_id, expires_at FROM sessions WHERE token = ?`, token).Scan(&userId, &expiresAt)
	if err == sql.ErrNoRows || err != nil {
		return ""
	}
	if time.Now().Unix() > expiresAt {
		return ""
	}
	return userId
}
