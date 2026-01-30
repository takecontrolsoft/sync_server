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

package impl

import (
	"bytes"
	"fmt"
	"image"
	"image/png"
	"io"
	"os"
	"path/filepath"

	ffmpeg "github.com/u2takey/ffmpeg-go"

	"github.com/disintegration/imaging"
	"github.com/takecontrolsoft/go_multi_log/logger"
	"github.com/takecontrolsoft/sync_server/server/config"
	"github.com/takecontrolsoft/sync_server/server/utils"
)

func GetFrameFromVideo(inFileName string, frameNum int) io.Reader {
	buf := bytes.NewBuffer(nil)
	err := ffmpeg.Input(inFileName).
		Filter("select", ffmpeg.Args{fmt.Sprintf("gte(n,%d)", frameNum)}).
		Output("pipe:", ffmpeg.KwArgs{"vframes": 1, "format": "image2", "vcodec": "mjpeg"}).
		WithOutput(buf, os.Stdout).
		Run()
	if err != nil {
		logger.Error(err)
	}
	return buf
}

func BuildVideoThumbnail(userName string, deviceId string, file string) (string, error) {
	// Normalize so Trash paths use forward slashes and ThumbnailBasePath puts thumb under Trash/Thumbnails.
	file = filepath.ToSlash(file)
	userDirName := filepath.Join(config.UploadDirectory, userName, deviceId)
	// Use ThumbnailBasePath so uploads-to-Trash and /img thumbnail lookup use the same path.
	thumbnailPath := ThumbnailBasePath(userDirName, file) + ".jpeg"
	filePath := filepath.Join(userDirName, file)

	reader := GetFrameFromVideo(filePath, 5)
	src, err := imaging.Decode(reader)
	if err != nil {
		return "", err
	}
	// Resize srcImage to width = 300px preserving the aspect ratio.
	resized := imaging.Resize(src, 300, 0, imaging.Lanczos)
	// Resize and crop the srcImage to fill the 250x250px area.
	thumbnail := imaging.Fill(resized, 250, 250, imaging.Center, imaging.Lanczos)

	// draw the srcImage over the backgroundImage at the (50, 50) position with opacity=0.5
	// playImage := imaging.OverlayCenter(thumbnail, image.Rect(0, 0, 50, 50), 0.5)
	err = os.MkdirAll(filepath.Dir(thumbnailPath), os.ModePerm)
	if err != nil {
		return "", err
	}
	err = imaging.Save(thumbnail, thumbnailPath)
	if err != nil {
		return "", err
	}

	return thumbnailPath, nil
}

// applyEXIFOrientation transforms the image according to EXIF Orientation (1-8).
// imaging: Rotate90 = 90° CCW, Rotate270 = 90° CW.
func applyEXIFOrientation(src image.Image, orientation int) image.Image {
	switch orientation {
	case 1:
		return src
	case 2:
		return imaging.FlipH(src)
	case 3:
		return imaging.Rotate180(src)
	case 4:
		return imaging.FlipV(src)
	case 5:
		return imaging.Rotate90(imaging.FlipH(src))
	case 6:
		return imaging.Rotate270(src)
	case 7:
		return imaging.Rotate270(imaging.FlipH(src))
	case 8:
		return imaging.Rotate90(src)
	default:
		return src
	}
}

func BuildImageThumbnail(userName string, deviceId string, file string) (string, error) {
	// Normalize so Trash paths use forward slashes and ThumbnailBasePath puts thumb under Trash/Thumbnails.
	file = filepath.ToSlash(file)
	userDirName := filepath.Join(config.UploadDirectory, userName, deviceId)
	// Use ThumbnailBasePath so uploads-to-Trash and /img thumbnail lookup use the same path.
	thumbnailPath := ThumbnailBasePath(userDirName, file)
	filePath := filepath.Join(userDirName, file)

	src, err := utils.GetImageFromFilePath(filePath)
	if err != nil {
		return "", err
	}

	// Apply EXIF orientation so thumbnail is displayed correctly (e.g. phone photos rotated 90°).
	metadataPath := MetadataPath(userDirName, file)
	orientation := GetOrientationFromMetadata(metadataPath)
	src = applyEXIFOrientation(src, orientation)

	// Resize srcImage to width = 300px preserving the aspect ratio.
	resized := imaging.Resize(src, 300, 0, imaging.Lanczos)
	// Resize and crop the srcImage to fill the 250x250px area.
	thumbnail := imaging.Fill(resized, 250, 250, imaging.Center, imaging.Lanczos)
	err = os.MkdirAll(filepath.Dir(thumbnailPath), os.ModePerm)
	if err != nil {
		return "", err
	}
	f, err := os.Create(thumbnailPath)
	if err != nil {
		return "", err
	}
	defer f.Close()
	png.Encode(f, thumbnail)
	return thumbnailPath, nil
}

func BuildAudioThumbnail(userName string, deviceId string, file string) (string, error) {
	return "", nil
}
