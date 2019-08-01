module github.com/navigaid/gitql

go 1.13

require (
	github.com/chzyer/readline v0.0.0-20180603132655-2972be24d48e
	github.com/jessevdk/go-flags v1.4.0
	github.com/libgit2/git2go v0.0.0-20190618093925-b2e2b2f71bb4
	github.com/mattn/go-runewidth v0.0.4 // indirect
	github.com/olekukonko/tablewriter v0.0.1
	github.com/pkg/errors v0.8.1
	golang.org/x/sys v0.0.0-20190726091711-fc99dfbffb4e // indirect
)

replace github.com/libgit2/git2go => github.com/navigaid/git2go v0.0.0-20190731180847-34ecd6adce7d
