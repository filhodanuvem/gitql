setup() {
    load 'test_helper/bats-support/load'
    load 'test_helper/bats-assert/load'
}

@test "Check exit code for select" {
  run ./gitql 'select * from commits'
  [ "$status" -eq 0 ]
}

@test "Check like for select" {
  run ./gitql -f json 'select message from commits where date > "2020-10-04" and date < "2020-10-06" and message like "update"'
  assert_output '[{"message":"update Dockerfile (#97)"}]'
}

@test "Check not like for select" {
  run ./gitql -f json 'select message from commits where date > "2020-10-04" and date < "2020-10-06" and message not like "update"'
  assert_output '[{"message":"Remove C dependencies (#95)"},{"message":"Build gitql using go build"},{"message":"Build gitql using go build"},{"message":"Build gitql using go build"},{"message":"Build gitql using go build"},{"message":"Build gitql using go build (#96)"},{"message":"Compile gitql with go-git support"}]'
}

@test "Check in for select" {
  run ./gitql -f json 'select distinct author from commits where "Tadeu" in author and date < "2021-01-01"'
  assert_output '[{"author":"Tadeu Zagallo"}]'
}

@test "Check count for select" {
  run ./gitql -f json 'select count(*) from commits where date > "2020-10-09" and date < "2020-10-17"'
  assert_output '[{"count":"17"}]'
}

@test "Select distinct should works" {
  run ./gitql -f json 'select distinct author from commits where date > "2020-10-09" and date < "2020-10-17" order by date asc'
  assert_output '[{"author":"Claudson Oliveira"},{"author":"Pedro Silva"},{"author":"Simon B"},{"author":"zweihander"}]'
}

@test "Select count should works" {
  run ./gitql -f json 'select count(*) from commits where date < "2018-01-05"'
  assert_output '[{"count":"192"}]'
}

@test "Select should works with order and limit" {
  run ./gitql -f json 'select date, message from commits where date < "2021-04-27" order by date desc limit 3'
  assert_output '[{"date":"2021-05-26 23:09:31","message":"Add smoke test about count"},{"date":"2021-05-26 23:06:42","message":"Smoke test on select discinct"},{"date":"2021-05-26 23:04:36","message":"Remove bats for windows"}]'
}

# bugs to be fixed

@test "Check incorrect usage of in for select" {
  skip "Should fail gracefully when using in the wrong way"
  run ./gitql -f json 'select distinct author from commits where author in "Tadeu" and date < "2021-01-01"'
  assert_output 'Unexpected T_IN after T_ID'
}

@test "Check incorrect json output of select" {
  skip "Should not return any other field than message"
  run ./gitql -f json 'select message from commits where date < "2021-05-27" order by date desc limit 3'
  assert_output '[{"message":"Add smoke test about count"},{"message":"Smoke test on select discinct"},{"message":"Remove bats for windows"}]'
}