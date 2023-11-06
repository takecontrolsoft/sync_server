package main

import (
	"os"
	"tcsi/sync_server/internal/errutil"
	"tcsi/sync_server/internal/serverapi"
)

const maxUploadFileSize = 5 * 1024 * 1024 * 1024 // 5 GB
const uploadPathVariable = "SYNC_STORAGE_PATH"
const portVariable = "SYNC_SERVER_PORT"

func main() {
	errutil.InitLogFile("server.log")
	serverapi.Configure(os.Getenv(uploadPathVariable), os.Getenv(portVariable), maxUploadFileSize)
}
