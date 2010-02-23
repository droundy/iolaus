package main;

import (
	"./git/git"
	"./git/plumbing"
	"./util/out"
	"./util/help"
)

func main() {
	git.AmInRepo("Must be in a repository to call whatsnew!")
	help.Init("see unrecorded changes.", plumbing.LsFiles)
	//plumbing.ReadTree(plumbing.Ref("HEAD"))
	out.Print(plumbing.DiffFiles([]string{}))
	for _,newf := range plumbing.LsOthers() {
		out.Print("Added ",newf)
	}
}
