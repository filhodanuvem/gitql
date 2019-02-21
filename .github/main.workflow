
workflow "New workflow" {
  on = "push"
  resolves = ["unit test"]
}

action "unit test" {
  uses = "docker://golang:latest"
  runs = "make test"
}

