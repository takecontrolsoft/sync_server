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
	"path/filepath"
	"strings"

	"github.com/disintegration/imaging"
	"github.com/takecontrolsoft/sync_server/server/utils"
)

// LooksLikeDocument returns true if the image appears to be a document, whiteboard,
// notebook, textbook or book page (flat, text-heavy). Uses a simple heuristic:
// high mean brightness and bimodal brightness (lots of light + dark pixels).
// May have false positives (e.g. white wall, bright sky) and false negatives
// (dark pages, low contrast). Only call for image files.
func LooksLikeDocument(fullPath string) bool {
	img, err := utils.GetImageFromFilePath(fullPath)
	if err != nil {
		return false
	}
	// Resize for speed; 200px width keeps enough detail
	small := imaging.Resize(img, 200, 0, imaging.Lanczos)
	bounds := small.Bounds()
	w, h := bounds.Dx(), bounds.Dy()
	if w < 10 || h < 10 {
		return false
	}
	var sum uint64
	light := 0   // brightness >= 240 (white/light background)
	dark := 0    // brightness <= 25 (text/lines)
	pixels := w * h
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, _ := small.At(x, y).RGBA()
			// 0-65535 -> 0-255
			bright := (r>>8 + g>>8 + b>>8) / 3
			if bright > 255 {
				bright = 255
			}
			sum += uint64(bright)
			if bright >= 240 {
				light++
			} else if bright <= 25 {
				dark++
			}
		}
	}
	mean := int(sum / uint64(pixels))
	// Document-like: mostly light background (mean high) and bimodal (light + dark pixels)
	// e.g. white page with black text, or notebook with lines. Thresholds tuned to catch more docs.
	lightDarkRatio := float64(light+dark) / float64(pixels)
	return mean >= 120 && lightDarkRatio >= 0.28
}

// IsImagePath returns true if the file extension is a common image type (case-insensitive).
func IsImagePath(path string) bool {
	ext := strings.ToLower(filepath.Ext(path))
	if ext == "" {
		return false
	}
	ext = ext[1:] // drop leading dot
	switch ext {
	case "jpg", "jpeg", "png", "gif", "bmp", "webp", "heic":
		return true
	}
	return false
}
