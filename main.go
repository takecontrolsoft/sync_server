package main

import (
	"internal/errors_util"
	"os"

	"takecontrolsoft.eu/sync/server/api"
)

const maxUploadFileSize = 5 * 1024 * 1024 * 1024 // 5 GB
const uploadPathVariable = "SYNC_STORAGE_PATH"
const portVariable = "SYNC_SERVER_PORT"

func main() {
	errors_util.InitLogFile("server.log")
	api.Configure(os.Getenv(uploadPathVariable), os.Getenv(portVariable), maxUploadFileSize)
}
