<img src="https://takecontrolsoft.eu/assets/img/takecontrolsoft-logo-green.png" alt="Sync Device by Take Control - software & infrastructure" width="25%">

[![Web site](https://img.shields.io/badge/Web_site-takecontrolsoft.eu-pink)](https://takecontrolsoft.eu/)
[![Linked in](https://img.shields.io/badge/Linked_In-takecontrolsoft-blue?style=flat&logo=linkedin)](https://www.linkedin.com/company/take-control-si/)
[![Go Doc](https://pkg.go.dev/badge/github.com/takecontrolsoft/sync_server.svg)](https://pkg.go.dev/github.com/takecontrolsoft/sync_server)
[![Docker Hub](https://img.shields.io/badge/Docker_Hub-takecontrolorg-blue?style=flat&logo=docker)](https://hub.docker.com/r/takecontrolorg/sync_server)

[![Project](https://img.shields.io/badge/Project-Sync_Device-darkred?style=flat&logo=github)](https://github.com/orgs/takecontrolsoft/projects/1)
[![License](https://img.shields.io/badge/License-Apache-purple)](https://www.apache.org/licenses/LICENSE-2.0)
[![Main](https://github.com/takecontrolsoft/sync_server/actions/workflows/main.yml/badge.svg)](https://github.com/takecontrolsoft/sync_server/actions/workflows/main.yml)

[![Release](https://img.shields.io/github/v/release/takecontrolsoft/sync_server.svg?style=flat&logo=github)](https://github.com/takecontrolsoft/sync_server/releases/latest)



# sync server
Golang server for uploading files and media files processing workflows.

## API endpoints

Base URL: `http://<host>:<port>` (e.g. `http://localhost:8080`).

| Method | Endpoint | Description |
|--------|----------|-------------|
| **POST** | `/upload` | Upload a file (multipart). Headers: `user` (JSON string), `date` (e.g. `2024-01`). Saves under `user/deviceId/` and creates thumbnails for images/videos. |
| **POST** | `/folders` | List folder structure (years and months) for a user and device. Body: `{ "User": "", "DeviceId": "" }`. Returns JSON array of `{ Year, Months[] }`. |
| **POST** | `/files` | List file paths in a folder. Body: `{ "UserData": { "User": "", "DeviceId": "" }, "Folder": "2024/01" }`. Returns JSON array of file path strings. Use `Folder: "Trash"` to list all files in Trash (paths like `Trash/2024/01/photo.jpg`). |
| **POST** | `/img` | Get image or thumbnail as PNG bytes. Body: `{ "UserData": { "User": "", "DeviceId": "" }, "File": "<path>", "Quality": "full" \| "" }`. Use `Quality: "full"` for original image; omit or empty for thumbnail. EXIF orientation is applied for correct display. |
| **GET** | `/stream` | Stream video/audio file with HTTP Range support (for playback/seek). Query: `User`, `DeviceId`, `File` (URL-encoded path, e.g. `2024/01/video.mp4`). |
| **POST** | `/move-to-trash` | Move files (and their thumbnails and metadata) to Trash. Body: `{ "UserData": { "User": "", "DeviceId": "" }, "Files": ["2024/01/photo.jpg", ...] }`. |
| **POST** | `/restore` | Restore files from Trash to their original folder (by path). Body: `{ "UserData": { "User": "", "DeviceId": "" }, "Files": ["Trash/2024/01/photo.jpg", ...] }`. |
| **POST** | `/regenerate-thumbnails` | Regenerate thumbnails for all media files (excluding Trash). Body: `{ "UserData": { "User": "", "DeviceId": "" } }`. Returns `{ "Regenerated": N }`. |
| **POST** | `/clean-orphan-thumbnails` | Delete thumbnail and metadata files that have no corresponding source file. Body: `{ "UserData": { "User": "", "DeviceId": "" } }`. Returns `{ "Removed": N }`. |
| **POST** | `/run-document-detection` | Run document detection (Python classifier if `SYNC_DOCUMENT_CLASSIFIER_PATH` is set, else built-in heuristic) on existing image files; move detected documents to Trash. Body: `{ "UserData": { "User": "", "DeviceId": "" } }`. Returns `{ "Moved": N }`. |
| **GET** | `/setup_info` | Placeholder; returns a short info message. |

# prerequisites
* MacOs
  
    `brew install exiftool`
  
    `export PATH=$PATH:/usr/local/bin`
  
    `brew install ffmpeg`
  
* Linux
  
    cd <your download directory>
    
    gzip -dc Image-ExifTool-12.96.tar.gz | tar -xf -cd Image-ExifTool-12.96
  
    sudo make install

    sudo apt install ffmpeg

* Windows

  The server looks for **exiftool** and **ffmpeg** in the same directory as `sync_server.exe`.  
  **Release artifacts** for Windows already include `exiftool.exe`, `ffmpeg.exe`, and `ffprobe.exe` next to the server executable.  
  For a custom build, place `exiftool.exe` and `ffmpeg.exe` in the folder where `sync_server.exe` is located, or install them system-wide and add their `bin` folder to `PATH`.

  - ExifTool: [exiftool.org](https://exiftool.org/) — download the Windows executable (e.g. `exiftool-13.33_64.zip`), extract and copy `exiftool(-k).exe` as `exiftool.exe` next to `sync_server.exe`.
  - FFmpeg: e.g. [BtbN/FFmpeg-Builds](https://github.com/BtbN/FFmpeg-Builds/releases) — copy `ffmpeg.exe` (and optionally `ffprobe.exe`) from the archive’s `bin` folder next to `sync_server.exe`.

## Auth (user names and passwords)

To use a **user ID** for folder paths (so usernames with invalid characters are safe and two users with the same display name don’t collide):

1. Set **`SYNC_AUTH_DB`** to the path of a SQLite file (e.g. `./sync_auth.db`). The server will create it and store users (id, username, password hash) and session tokens there.
2. Set **`SYNC_ADMIN_USER`** and **`SYNC_ADMIN_PASSWORD`** to bootstrap the first user (used only when the DB has no users).
3. **POST /auth/login** – body `{ "User": "<username/email>", "Password": "" }` → returns `{ "Token": "...", "UserId": "<uuid>" }`. Use **UserId** (not the username) in upload, /folders, /files, /img, /stream so paths on disk are `UserId/DeviceId/...`.
4. **POST /auth/register** – body `{ "User": "", "Password": "" }` creates a new user and returns `{ "Token": "...", "UserId": "..." }`.
5. **Folder layout on disk**: when auth is enabled, files are under `UserId/DeviceId/year/month/` (UserId is a UUID from the DB; username is only stored in the DB).

If `SYNC_AUTH_DB` is not set, the **User** value from the client is used as the folder name (legacy behaviour).

## Optional: document-to-Trash detection

Set **`SYNC_DOCUMENT_TO_TRASH=1`** (or `true` / `yes`) so that uploaded **images** that look like documents (whiteboard, notebook, textbook, book page) are automatically moved to Trash. The server uses a simple heuristic: high mean brightness and many light + dark pixels (typical for text on white background). This can have false positives (e.g. bright sky, white wall) and false negatives (dark pages). Disable the option if too many normal photos are moved.

# Building and release
Go to CONTRIBUTING.md for more instructions.
