#!/usr/bin/env bats

@test "Check version" {
  result="$(./gitql -v)"
  [ "$result" == "Gitql 2.2.0" ]
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