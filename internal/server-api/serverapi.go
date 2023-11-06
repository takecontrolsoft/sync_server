// Package serverapi provides API for files synchronization
package serverapi

import (
	"fmt"
	"log"
	"net/http"
)

var uploadDirectory string
var maxUploadFileSize int64

// Initial entry point for configuring and starting th server
func Configure(uploadDirectoryParam string, portParam string, maxUploadFileSizeParam int64) {
	uploadDirectory = uploadDirectoryParam
	maxUploadFileSize = maxUploadFileSizeParam

	http.HandleFunc("/uploadFile", uploadFileHandler())

	fs := http.FileServer(http.Dir(uploadDirectory))
	http.Handle("/files/", http.StripPrefix("/files", fs))

	log.Printf("Server started on localhost port: %s", portParam)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", portParam), nil))
}
