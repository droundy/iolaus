package main;

import (
	"os"
	"goopt"
	"./git/git"
	"./git/plumbing"
	"./util/out"
	"./util/error"
	"./util/help"
	"./util/cook"
)

var shortlog = goopt.String([]string{"-m","--patch"}, "COMMITNAME",
	"name of commit")

func main() {
	git.AmInRepo("Must be in a repository to call record!")
	help.Init("record changes.", plumbing.LsFiles)
	plumbing.ReadTree(plumbing.Ref("HEAD"))

	for _,f := range plumbing.DiffFilesModified([]string{}) {
		out.Print("Considering changes to",f)
		plumbing.UpdateIndex(f)
	}
	for _,newf := range plumbing.LsOthers() {
		out.Print("Considering adding",newf)
		plumbing.UpdateIndex(newf)
	}
	if *shortlog != "COMMITNAME" {
		c := plumbing.CommitTree(plumbing.WriteTree(),
			[]plumbing.Commitish{plumbing.Ref("HEAD")}, *shortlog)
		plumbing.UpdateRef("HEAD", c)
	} else {
		cook.SetRaw()
		defer cook.SetCooked()
		x := make([]byte,1)
		os.Stdin.Read(x)
		out.Print("byte is",string(x))
		error.FailOn(os.NewError("FIXME, need to read log"))
	}
}
