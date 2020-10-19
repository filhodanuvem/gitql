#!/usr/bin/env bats

setup() {
  export branch=$(git branch | grep "\*" | cut -d " " -f 2)
}

teardown() {
  git checkout $branch &> /dev/null 
}

@test "Check switching to existing branch" {
  run ./gitql 'use master'
  [ "$status" -eq 0 ]
}

@test "Check switching to nonexistent branch" {
  run ./gitql 'use this-is-not-a-branch'
  [ "$status" -eq 1 ]
}
