<img src="https://takecontrolsoft.eu/wp-content/uploads/2023/11/TakeControlTransparentGreenLogo-1.png" alt="Sync Device by Take Control - software & infrastructure" width="25%">

[![Web site](https://img.shields.io/badge/Web_site-takecontrolsoft.eu-pink)](https://takecontrolsoft.eu/)
[![Linked in](https://img.shields.io/badge/Linked_In-page-blue)](https://www.linkedin.com/company/take-control-si/)
[![Docker Hub](https://img.shields.io/badge/Docker_Hub-repo-blue)](https://hub.docker.com/repository/docker/takecontrolorg/sync_server/general)

[![Project](https://img.shields.io/badge/Project-Sync_Device-darkred)](https://github.com/orgs/takecontrolsoft/projects/1)
[![License](https://img.shields.io/badge/License-Apache-purple)](https://www.apache.org/licenses/LICENSE-2.0)
[![Main](https://github.com/takecontrolsoft/sync_server/actions/workflows/main.yml/badge.svg)](https://github.com/takecontrolsoft/sync_server/actions/workflows/main.yml)
[![Pull Request](https://github.com/takecontrolsoft/sync_server/actions/workflows/pull_request.yml/badge.svg)](https://github.com/takecontrolsoft/sync_server/actions/workflows/pull_request.yml)

[![Release](https://img.shields.io/github/v/release/takecontrolsoft/sync_server.svg)](https://github.com/takecontrolsoft/sync_server/releases/latest)

<!-- ![GitHub release (by tag)](https://img.shields.io/github/downloads/takecontrolsoft/sync_server/v0.0.1-alpha/total)
![Docker Pulls](https://img.shields.io/docker/pulls/takecontrolorg/sync_server) -->

# sync server
Golang server for uploading files and media files processing workflows.

go get github.com/takecontrolsoft/go_multi_log@v1.0.1

go build -v ./...

go build -o bin

bin/sync_server.exe /help

bin/sync_server.exe -p 3000 -d C:\Users\desis\Pictures\FileSyncTest\ -l C:\Users\desis\Pictures\FileSyncTest\ -n 5
http://localhost:3000/files


bin/sync_server.exe -p 3000 -d /photos/ -l /log/ -n 5

 godoc -http=:8081 -index
 http://localhost:8081/pkg/


