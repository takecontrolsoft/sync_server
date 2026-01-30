<img src="https://takecontrolsoft.eu/assets/img/takecontrolsoft-logo-green.png" alt="Sync Device by Take Control - software & infrastructure" width="25%">

## 1.0.8 Release notes (2026-01-29)

### Fixes
* Streaming: use `io.Copy` instead of `WriteTo` for full-file responses (fixes build on Go 1.22+)
* CI: Docker workflow no longer uploads bake-meta/digest artifacts; digests passed via job outputs; clean-artifacts job removed

## 1.0.7 Release notes (2026-01-29)

### Enhancements
* Auth: user folders on disk use lowercase email; login/register return UserId for paths
* Admin APIs: POST `/regenerate-thumbnails`, POST `/clean-orphan-thumbnails`, POST `/run-document-detection`
* Document detection: fallback to built-in heuristic when Python classifier is unset or fails; relaxed thresholds
* Standalone build workflow (`.github/workflows/build.yml`) for CI
* Tests for admin endpoints (regenerate-thumbnails, clean-orphan-thumbnails, run-document-detection)

### Fixes
* Document detection returns `{"Moved": 0}` (200) when user directory is empty instead of 500
* `IsImagePath` is case-insensitive for extensions (e.g. .JPG)

## 1.0.6 Release notes (2025-07-27)

### Enhancements
* Added support for returning full-size images
### Fixes
* Fixed issue where metadata files would fail to write when parent directories didn't exist - directories are now created automatically

## 1.0.5 Release notes (2024-09-03)

### Enhancements
* Generating Video thumbnails.
* Improved Photos thumbnails.
* Extracting media file properties to json metadata fails.

## 1.0.4 Release notes (2024-08-21)

### Enhancements
* Creating thumbnails of images while uploading.

## 1.0.3 Release notes (2024-08-20)

### Enhancements
* Implemented API `/img` that returns thumbnails streams.
* Implemented API `/delete-all` that deletes all the files for the device and user from the server.


## 1.0.2 Release notes (2024-08-18)

### Enhancements
* Implemented API `/folders` that returns all the folders in a tree for the device and user.
* Implemented API `/files` that returns all the files under a specific folder for the device and user.

## 1.0.1 Release notes (2024-08-02)

### Enhancements
* Files stored in monthly folders.

## 1.0.0 Release notes (2024-08-01)

### Enhancements
* Store files in folders by username, deviceId and date of the files.
* Improve logging. Log to a file.

## 0.0.1-alpha Release notes (2024-01-09)

### Enhancements
* The initial alpha version of the Sync Server for media files.

### Compatibility
* Linux, Windows and Mac
* Go lang 1.21
