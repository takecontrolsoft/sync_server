// Package serverapi provides API for files synchronization
package api

import (
	"fmt"
	"log"
	"net/http"

	"takecontrolsoft.eu/sync/server/config"
	"takecontrolsoft.eu/sync/server/impl"
)

// Initial entry point for configuring and starting th server
func Configure(uploadDirectoryParam string, portParam string, maxUploadFileSizeParam int64) {
	config.UploadDirectory = uploadDirectoryParam
	config.MaxUploadFileSize = maxUploadFileSizeParam

	http.HandleFunc("/upload", impl.UploadHandler())

	fs := http.FileServer(http.Dir(config.UploadDirectory))
	http.Handle("/files/", http.StripPrefix("/files", fs))

	log.Printf("Server started on localhost port: %s", portParam)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", portParam), nil))
}
