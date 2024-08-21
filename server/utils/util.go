package utils

import (
	"bytes"
	"encoding/json"
	"image"
	"image/color"

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

func GetImageFromFilePath(filePath string) (image.Image, error) {
	reader, err := os.Open(filePath)
	if err != nil {
		logger.Error(err)
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
