
workflow "New workflow" {
  on = "push"
  resolves = ["unit test"]
}

action "unit test" {
  uses = "actions/docker/cli@8cdf801b322af5f369e00d85e9cf3a7122f49108"
  runs = "ls"
}

