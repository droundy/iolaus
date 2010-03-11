package main;

import (
	"./git/git"
	"./git/color"
	"./git/plumbing"
	"./util/out"
	"./util/help"
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
		out.Print(color.String("Added "+newf, color.Meta))
	}
}
