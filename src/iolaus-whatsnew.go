package main;

import (
	"goopt"
	"./git/git"
	"./git/plumbing"
	"./util/help"
	"./util/exit"
	"./util/error"
	"./iolaus/core"
	"./iolaus/prompt"
)

var description = func() string {
	return `
whatsnew gives you a view of what changes you've made in your working
copy that haven't yet been recorded.
`}

func main() {
	goopt.Vars["Verb"] = "Display"
	goopt.Vars["verb"] = "display"
	*prompt.All = true
	help.Init("see unrecorded changes.", description, plumbing.LsFiles)
	git.AmInRepo("Must be in a repository to call whatsnew!")
	//plumbing.ReadTree(git.Ref("HEAD"))
	ds,e := core.DiffFiles([]string{})
	error.FailOn(e)
	prompt.Run(ds, func (f core.FileDiff) {
		f.Print()
	})
	exit.Exit(0)
}
