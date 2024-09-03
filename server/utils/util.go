/* Copyright 2024 Take Control - Software & Infrastructure

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

package utils

import (
	"bufio"
	"bytes"
	"encoding/json"
	"image"
	"image/color"
	"path/filepath"
	"runtime"

	"image/draw"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/go-errors/errors"
	"github.com/takecontrolsoft/go_multi_log/logger"
	"github.com/takecontrolsoft/sync_server/server/mediatypes"
)

func RenderIfError(err error, w http.ResponseWriter, statusCode int) bool {
	if err != nil {
		http.Error(w, err.Error(), statusCode)
		logger.Error(err)
		return true
	}
	return false
}

func RenderError(w http.ResponseWriter, err error, statusCode int) {
	w.WriteHeader(statusCode)
	w.Write([]byte(err.Error()))
	logger.Error(err)
}

func IsAllowedFileType(fileType string, w http.ResponseWriter) bool {
	allowed := []string{"image/", "video/", "audio/"}
	result := false
	for i := range allowed {
		result = result || strings.HasPrefix(fileType, allowed[i])
	}
	return result
}

func GetMediaType(fileType string) mediatypes.MediaType {
	allowed := []string{"image/", "video/", "audio/"}
	result := ""
	for i := range allowed {
		if strings.HasPrefix(fileType, allowed[i]) {
			result = allowed[i]
		}
	}
	switch result {
	case "image/":
		return mediatypes.Image
	case "video/":
		return mediatypes.Video
	case "audio/":
		return mediatypes.Audio
	default:
		return mediatypes.Unknown
	}
}

func GenerateRandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	seed := rand.NewSource(time.Now().UnixNano())
	random := rand.New(seed)

	result := make([]byte, length)
	for i := range result {
		result[i] = charset[random.Intn(len(charset))]
	}
	return string(result)
}

func JsonReaderFactory(in interface{}) (io.Reader, error) {
	buf := bytes.NewBuffer(nil)
	enc := json.NewEncoder(buf)
	err := enc.Encode(in)
	if err != nil {
		return nil, errors.Errorf("creating reader: error encoding data: %s", err)
	}
	return buf, nil
}

func ResizeImage(img *image.RGBA, height int) *image.RGBA {
	if height < 50 {
		return img
	}
	bounds := img.Bounds()
	imgHeight := bounds.Dy()
	if height >= imgHeight {
		return img
	}
	imgWidth := bounds.Dx()
	resizeFactor := float32(imgHeight) / float32(height)
	ratio := float32(imgWidth) / float32(imgHeight)
	width := int(float32(height) * ratio)
	resizedImage := image.NewRGBA(image.Rect(0, 0, width, height))
	var imgX, imgY int
	var imgColor color.Color
	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			imgX = int(resizeFactor*float32(x) + 0.5)
			imgY = int(resizeFactor*float32(y) + 0.5)
			imgColor = img.At(imgX, imgY)
			resizedImage.Set(x, y, imgColor)
		}
	}
	return resizedImage
}

func ImageToRGBA(src image.Image) *image.RGBA {

	if dst, ok := src.(*image.RGBA); ok {
		return dst
	}

	b := src.Bounds()
	dst := image.NewRGBA(image.Rect(0, 0, b.Dx(), b.Dy()))
	draw.Draw(dst, dst.Bounds(), src, b.Min, draw.Src)
	return dst
}

func GetThumbnailFileAddedExtension(filePath string) (string, error) {
	reader, err := os.Open(filePath)
	if err != nil {
		logger.Error(err)
		return "", err
	}
	defer reader.Close()
	b := bufio.NewReader(reader)
	n, _ := b.Peek(512)
	fileType := http.DetectContentType(n)
	mediaType := GetMediaType(fileType)
	if mediaType != mediatypes.Image {
		return ".jpeg", nil
	}
	return "", nil
}

func GetImageFromFilePath(filePath string) (image.Image, error) {
	reader, err := os.Open(filePath)
	if err != nil {
		logger.Error(err)
		return nil, err
	}
	defer reader.Close()

	reader.Seek(0, 0)

	m, _, err := image.Decode(reader)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	return m, err
}

func GetExecutablePath() (string, error) {
	var exPath string
	ex, err := os.Executable()
	if err != nil {
		return "", err
	}
	exPath = filepath.Dir(ex)
	return exPath, nil
}

func GetToolPath(toolName string) (string, error) {
	executablePath := ""
	extension := ""
	exifToolFile := toolName
	if runtime.GOOS == "windows" {
		extension = ".exe"
		path, err := GetExecutablePath()
		if err != nil {
			return "", err
		}
		executablePath = path
		exifToolFile = filepath.Join(executablePath, toolName, extension)
	}
	_, err := os.Stat(exifToolFile)
	if err != nil {
		return "", err
	}
	return exifToolFile, nil
}
