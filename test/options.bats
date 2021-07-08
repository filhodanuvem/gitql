#!/usr/bin/env bats

# test that will fail outside of pipeline because
# it requires an override of the version.txt
# you can check how it's done on github actions files.
@test "Check version with -v" {
  result="$(./gitql -v)"
  [ "$result" != "Gitql latest" ]
}

@test "Check version" {
  result="$(./gitql version)"
  [ "$result" != "Gitql latest" ]
}

@test "Check table commits on -s" {
  result="$(./gitql -s | grep commits)"
  [ "$result" == "commits" ] 
}

@test "Check table refs on -s" {
  result="$(./gitql -s | grep refs)"
  [ "$result" == "refs" ] 
}

@test "Check table tags on -s" {
  result="$(./gitql -s | grep tags)"
  [ "$result" == "tags" ] 
}

@test "Check table branches on -s" {
  result="$(./gitql -s | grep branches)"
  [ "$result" == "branches" ] 
}

@test "Check exit code for help" {
  result="$(./gitql)"
  [ "$?" == "0" ] 
}