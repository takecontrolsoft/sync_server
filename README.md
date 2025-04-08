<img src="https://takecontrolsoft.eu/assets/img/takecontrolsoft-logo-green.png" alt="Sync Device by Take Control - software & infrastructure" width="25%">

[![Web site](https://img.shields.io/badge/Web_site-takecontrolsoft.eu-pink)](https://takecontrolsoft.eu/)
[![Linked in](https://img.shields.io/badge/Linked_In-takecontrolsoft-blue?style=flat&logo=linkedin)](https://www.linkedin.com/company/take-control-si/)
[![Go Doc](https://pkg.go.dev/badge/github.com/takecontrolsoft/sync_server.svg)](https://pkg.go.dev/github.com/takecontrolsoft/sync_server)
[![Docker Hub](https://img.shields.io/badge/Docker_Hub-takecontrolorg-blue?style=flat&logo=docker)](https://hub.docker.com/r/takecontrolorg/sync_server)

[![Project](https://img.shields.io/badge/Project-Sync_Device-darkred?style=flat&logo=github)](https://github.com/orgs/takecontrolsoft/projects/1)
[![License](https://img.shields.io/badge/License-Apache-purple)](https://www.apache.org/licenses/LICENSE-2.0)
[![Main](https://github.com/takecontrolsoft/sync_server/actions/workflows/main.yml/badge.svg)](https://github.com/takecontrolsoft/sync_server/actions/workflows/main.yml)

[![Release](https://img.shields.io/github/v/release/takecontrolsoft/sync_server.svg?style=flat&logo=github)](https://github.com/takecontrolsoft/sync_server/releases/latest)



# sync server
Golang server for uploading files and media files processing workflows.

# prerequisites
* MacOs
  
    `brew install exiftool`
  
    `export PATH=$PATH:/usr/local/bin`
  
    `brew install ffmpeg`
  
* Linux
  
    cd <your download directory>
    
    gzip -dc Image-ExifTool-12.96.tar.gz | tar -xf -cd Image-ExifTool-12.96
  
    sudo make install

    sudo apt install ffmpeg





# Building and release
Go to CONTRIBUTING.md for more instructions.
