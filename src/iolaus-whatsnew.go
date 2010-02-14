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
	out.Print(plumbing.Diff([]string{}))
	for _,newf := range plumbing.LsOthers() {
		out.Print("Added",newf)
	}
}
