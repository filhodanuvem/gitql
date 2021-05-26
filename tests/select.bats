#!/bin/bash

@test "Check exit code for select" {
  result="$(./gitql 'select * from commits')"
  [ "$?" == "0" ] 
}

@test "Select discting should work" {
  run ./gitql 'select distinct author from commits'
  [ "$status" -eq 0 ]
}

@test "Select count should work" {
  run ./gitql 'select count(*) from commits'
  [ "$status" -eq 0 ]
}@test "Select discting should work" {
  run ./gitql 'select distinct author from commits'
  [ "$status" -eq 0 ]
}

@test "Select count should work" {
  run ./gitql 'select count(*) from commits'
  [ "$status" -eq 0 ]
}