package main;

import (
	"./git/git"
	"./git/plumbing"
	"./util/error"
	"./util/help"
)

func main() {
	git.AmInRepo("Must be in a repository to call whatsnew!")
	help.Init("see unrecorded changes.", plumbing.LsFiles)
	p := plumbing.Diff([]string{})
	error.Print(p)
}
