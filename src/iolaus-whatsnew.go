package main;

import (
	"./git/git"
	"./git/color"
	"./git/plumbing"
	"./util/out"
	"./util/help"
	"./util/exit"
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
	p := plumbing.DiffFiles([]string{}).String()
	if p != "" {
		out.Print(p)
	}
	for _,newf := range plumbing.LsOthers() {
		out.Print(color.String("Added "+newf, color.Meta))
	}
	exit.Exit(0)
}
