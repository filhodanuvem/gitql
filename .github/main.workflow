workflow "Functional tests" {
  on = "push"
  resolves = ["build"]
}

action "build" {
  uses = "./.github/functional-tests/"
}