package main;

import (
	"./git/git"
	"./git/plumbing"
	"./util/out"
	"./util/help"
	"./util/exit"
)

func main() {
	git.AmInRepo("Must be in a repository to call whatsnew!")
	help.Init("see unrecorded changes.", plumbing.LsFiles)
	//plumbing.ReadTree(plumbing.Ref("HEAD"))
	p := plumbing.DiffFiles([]string{}).String()
	if p != "" {
		out.Print(p)
	}
	for _,newf := range plumbing.LsOthers() {
		out.Print("Added ",newf)
	}
	exit.Exit(0)
}
