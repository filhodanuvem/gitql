#!/bin/bash

@test "Check exit code for select" {
  result="$(./gitql 'select * from commits')"
  [ "$?" == "0" ] 
}

@test "Select discting should work" {
  run ./gitql 'select distinct author from commits'
  [ "$status" -eq 0 ]
}