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

# How to contribute

## Get packages
`go get github.com/takecontrolsoft/go_multi_log@v1.0.1`

## Build go server
`go build -v ./...`

or

`go build -o bin`

# How to run sync server
## Open sync server help
`bin/sync_server.exe /help`

## Example of server parameters
`bin/sync_server.exe -p 3000 -d C:\Users\{username}\Pictures\FileSyncTest\ -l C:\Users\{username}\Pictures\FileSyncTest\ -n 5`

or

`bin/sync_server.exe -p 3000 -d /photos/ -l /log/ -n 5`

## To browse server
http://localhost:3000/files

# Sync server documentation
## To build documentation
`godoc -http=:8081 -index`
## Brows documentation
 http://localhost:8081/pkg/


# Docker image
## To build an image named "tc" run:
`docker build . -t tc -f Dockerfile.linux`

## To list docker images run:
`docker images`

## To delete docker image "tc" run:
`docker rmi tc:latest -f`

## To run docker image "tc" run:
`docker run --name t1 -p 3000:3000 tc:latest -e "LOG_LEVEL=3" -v /photos:./bin /logs:./bin`

