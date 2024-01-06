variable "TAG" {
  default = "latest"
}
  
group "default" {
  targets = ["sync_server"]
}

target "sync_server" {
  dockerfile = "Dockerfile"
  tags = ["docker.io/takecontrolorg/sync_server:${TAG}"]
  platforms = ["linux/amd64", "linux/arm64", "darwin/amd64", "darwin/arm64", "windows/amd64", "windows/arm64"]
}