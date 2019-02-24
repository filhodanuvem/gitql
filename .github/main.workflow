workflow "New workflow" {
  on = "push"
  resolves = ["build"]
}

action "build" {
  uses = "./.github/functional-tests/"
}
