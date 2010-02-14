package main;

import (
	"./git/git"
	"./git/plumbing"
	"./util/error"
	"./util/help"
)

func main() {
	git.AmInRepo("Must be in a repository to call whatsnew!")
	help.Init("see unrecorded changes.")
	fs,_ := plumbing.LsFiles()
	error.Print(fs)
	p := plumbing.DiffFiles([]string{})
	error.Print(p)
}
