group "default" {
  targets = ["linux"]
}

target "linux" {
  dockerfile = "Dockerfile.linux"
  tags = ["docker.io/takecontrolorg/sync_server_linux"]
  platforms = ["linux/amd64", "linux/arm64"]
}