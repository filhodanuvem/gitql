#!/bin/bash

@test "Check exit code for select" {
  result="$(./gitql 'select * from commits')"
  [ "$?" == "0" ] 
}