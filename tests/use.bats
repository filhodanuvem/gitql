#!/usr/bin/env bats

@test "Check switching to existing branch" {
  result="$(./gitql 'use master')"
  [ "$?" == "0" ]
}

@test "Check switching to nonexistent branch" {
  result="$(./gitql 'use this-is-not-a-branch')"
  [ "$?" == "1" ]
}
