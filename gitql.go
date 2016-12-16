package main

import "os"

func main() {
	cmd := new(Gitql)
	os.Exit(cmd.Run())
}
