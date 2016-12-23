package main

import "os"

func main() {
	os.Exit(new(Gitql).Run())
}
