variable "TAG" {
  default = "latest"
}
  
group "default" {
  targets = ["sync_server"]
}

target "sync_server" {
  tags = ["docker.io/takecontrolorg/sync_server:${TAG}"]
}

target "sync_server_linux" {
  inherits = ["sync_server"]
  dockerfile = "Dockerfile.linux"
  architectures = ["amd64", "arm64"]
}

