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
	"bufio"
	"bytes"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/takecontrolsoft/go_multi_log/logger"
	"github.com/takecontrolsoft/sync_server/server/config"
)

// RunDocumentClassifierAsync runs the configured Python/exe classifier in a new
// goroutine. After sync, if stdout contains "document", the file is moved to Trash.
// fullPath is the absolute path to the image; userDir and relPath are for MoveRelativePathToTrash.
func RunDocumentClassifierAsync(fullPath, userDir, relPath string) {
	if config.DocumentClassifierPath == "" {
		return
	}
	go func() {
		RunDocumentClassifierSync(fullPath, userDir, relPath)
	}()
}

// RunDocumentClassifierSync runs the classifier and, if stdout contains "document", moves the file to Trash.
// Call this after thumbnail creation so the thumbnail is moved to Trash/Thumbnails together with the file.
func RunDocumentClassifierSync(fullPath, userDir, relPath string) {
	if config.DocumentClassifierPath == "" {
		return
	}
	out, err := runClassifier(fullPath)
	if err != nil {
		logger.ErrorF("Document classifier failed for %s: %v", fullPath, err)
		return
	}
	if strings.Contains(strings.ToLower(string(out)), "document") {
		MoveRelativePathToTrash(userDir, relPath)
	}
}

// runClassifier runs the script or exe with image path as single arg; returns stdout.
func runClassifier(imagePath string) ([]byte, error) {
	path := strings.TrimSpace(config.DocumentClassifierPath)
	if path == "" {
		return nil, nil
	}
	absPath := path
	if !filepath.IsAbs(path) {
		// Resolve relative to BinDirectory (next to sync_server) so scripts/ lives there
		absPath = filepath.Join(config.BinDirectory, path)
	}
	ext := strings.ToLower(filepath.Ext(absPath))
	var cmd *exec.Cmd
	if ext == ".py" {
		cmd = exec.Command("python", absPath, imagePath)
	} else if ext == ".exe" || ext == "" {
		cmd = exec.Command(absPath, imagePath)
	} else {
		cmd = exec.Command(absPath, imagePath)
	}
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		if stderr.Len() > 0 {
			logger.ErrorF("Classifier stderr: %s", bufio.NewScanner(strings.NewReader(stderr.String())))
		}
		return nil, err
	}
	return stdout.Bytes(), nil
}
