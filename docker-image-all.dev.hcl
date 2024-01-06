group "default" {
  targets = ["linux"]
}

target "linux" {
  dockerfile = "Dockerfile.linux"
  tags = ["docker.io/takecontrolorg/sync_server_linux"]
}