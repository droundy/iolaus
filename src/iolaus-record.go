package main;

import (
	"./git/git"
	"./git/plumbing"
	"./util/out"
	"./util/help"
)

func main() {
	git.AmInRepo("Must be in a repository to call record!")
	help.Init("record changes.", plumbing.LsFiles)
	plumbing.ReadTree(plumbing.Ref("HEAD"))

	for _,f := range plumbing.DiffFilesModified([]string{}) {
		out.Print("Considering changes to",f)
	}
	for _,newf := range plumbing.LsOthers() {
		out.Print("Considering adding",newf)
	}
}
