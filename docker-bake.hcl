variable "TAG" {
  default = "latest"
}
  
group "default" {
  targets = ["sync_server_linux", "sync_server_macos", "sync_server_windows"]
}

target "sync_server" {
  tags = ["docker.io/takecontrolorg/sync_server:${TAG}"]
}

target "sync_server_linux" {
  inherits = ["sync_server"]
  dockerfile = "Dockerfile.linux"
  platforms = ["amd64", "arm64"]
}

target "sync_server_macos" {
  inherits = ["sync_server"]
  dockerfile = "Dockerfile.macos"
  platforms = ["amd64", "arm64"]
}

target "sync_server_windows" {
  inherits = ["sync_server"]
  dockerfile = "Dockerfile.windows"
  platforms = ["amd64", "arm64"]
}

