variable "TAG" {
  default = "latest"
}
  
group "default" {
  targets = ["sync_server.linux", "sync_server.windows"]
}

target "sync_server" {
   tags = ["docker.io/takecontrolorg/sync_server:${TAG}"]
}

target "sync_server.linux" {
  inherits = ["sync_server"]
  dockerfile = "Dockerfile.linux"
  platforms = ["linux/amd64", "linux/arm64"]
}

target "sync_server.windows" {
  inherits = ["sync_server"]
  dockerfile = "Dockerfile.windows"
  platforms = ["windows/amd64", "windows/arm64"]
}