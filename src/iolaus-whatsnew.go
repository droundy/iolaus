package main;

import (
	"./git/git"
	"./git/plumbing"
	"./util/help"
	"./util/exit"
	"./util/error"
	"./iolaus/core"
)

var description = func() string {
	return `
whatsnew gives you a view of what changes you've made in your working
copy that haven't yet been recorded.
`}

func main() {
	help.Init("see unrecorded changes.", description, plumbing.LsFiles)
	git.AmInRepo("Must be in a repository to call whatsnew!")
	//plumbing.ReadTree(git.Ref("HEAD"))
	ds,e := core.DiffFiles([]string{})
	error.FailOn(e)
	// The following should all be factored apart next...
	for _,d := range ds {
		d.Print()
	}
	exit.Exit(0)
}
