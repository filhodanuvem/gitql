workflow "New workflow" {
  on = "push"
  resolves = ["build"]
}

action "build" {
  uses = "docker://golang:latest"
  runs = "go"
  args = " get -u -d github.com/cloudson/gitql"
}
