variable "TAG" {
  default = "latest"
}
  
group "default" {
  targets = ["sync_server"]
}

target "sync_server" {
  tags = ["docker.io/takecontrolorg/sync_server:${TAG}"]
  dockerfile = "Dockerfile"
  platforms = ["linux/amd64", "linux/arm64"]
}

