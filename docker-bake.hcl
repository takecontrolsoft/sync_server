variable "TAG" {
  default = "latest"
}
  
group "default" {
  targets = ["sync_server.linux", "sync_server.windows"]
}

target "sync_server.linux" {
  dockerfile = "Dockerfile.linux"
  tags = ["docker.io/takecontrolorg/sync_server:${TAG}"]
  platforms = ["linux/amd64", "linux/arm64", "darwin/amd64", "darwin/arm64", "windows/amd64", "windows/arm64"]
}

target "sync_server.windows" {
  dockerfile = "Dockerfile.windows"
  tags = ["docker.io/takecontrolorg/sync_server:${TAG}"]
  platforms = ["windows/amd64", "windows/arm64"]
}