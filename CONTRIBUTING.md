
# How to contribute

## Get packages
`go get github.com/takecontrolsoft/go_multi_log@v1.0.1`

## Build go server
* to local folder - `go build -v ./...`

* to bin folder `go build -o bin/`

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
```bash
go get golang.org/x/tools/cmd/godoc
export GOPATH=$HOME/go 
export GOROOT=/usr/local/go/bin
export PATH=$PATH:$GOPATH/bin
export PATH=$PATH:$GOROOT/bin
mkdir -p $GOPATH $GOPATH/src $GOPATH/pkg $GOPATH/bin
go install golang.org/x/tools/cmd/godoc@latest
godoc -http=:8081 -index
```
## Brows documentation
 http://localhost:8081/pkg/


# Docker image
## To build an image named "tc" run:
`docker build . -t tc -f Dockerfile.linux --platform linux/amd64`

## To list docker images run:
`docker images`

## To delete docker image "tc" run:
`docker rmi tc:latest -f`

## To run docker image "tc" run:
`sudo docker run --name mobisync -p 3000:3000 --mount type=bind,source=/mobisync,target=/data -e "LOG_LEVEL=3" takecontrolorg/sync_server:latest --add-host host.docker.internal:host-gateway`

## How to release
`git tag v1.0.0`      
`git push --tags`   

## How to update package by commit hash
`go get -u "github.com/takecontrolsoft/go_multi_log@d020e35eaecbfb8bfe32d368439d926f58d06d30"`