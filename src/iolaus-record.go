package main;

import (
	"os"
	"bufio"
	"goopt"
	git "./git/git"
	"./git/plumbing"
	"./util/out"
	"./util/error"
	"./util/help"
	"./util/cook"
	"./iolaus/test"
	hashes "./gotgo/slice(git.Commitish)"
)

var shortlog = goopt.String([]string{"-m","--patch"}, "COMMITNAME",
	"name of commit")
var all = goopt.Flag([]string{"-a","--all"}, []string{"--interactive"},
	"record all patches", "prompt for patches interactively")

var description = func() string {
	return `
Record is used to name a set of changes and record the patch to the
repository.
`}

func main() {
	help.Init("record changes.", description, plumbing.LsFiles)
	git.AmInRepo("Must be in a repository to call record!")
	//plumbing.ReadTree(plumbing.Ref("HEAD"))

	if *all {
		for _,f := range plumbing.DiffFilesModified([]string{}) {
			out.Print("Considering changes to ",f)
			plumbing.UpdateIndex(f)
		}
		for _,newf := range plumbing.LsOthers() {
			out.Print("Considering adding ",newf)
			plumbing.UpdateIndex(newf)
		}
	} else {
		unraw := cook.SetRaw()
		defer cook.Undo(unraw)
		for _,f := range plumbing.DiffFilesModified([]string{}) {
			c,e := out.PromptForChar("Record changes to %s? ", f)
			switch c {
			case 'q','Q': error.Exit(e)
			case 'y','Y':
				out.Print("Dealing with file ",f)
				plumbing.UpdateIndex(f)
			case 'n','N': out.Print("Ignoring changes to file ",f)
			}
		}
		for _,f := range plumbing.LsOthers() {
			c,e := out.PromptForChar("Record addition of %s? ", f)
			switch c {
			case 'q','Q': error.Exit(e)
			case 'y','Y':
				out.Print("Adding file ",f)
				plumbing.UpdateIndex(f)
			case 'n','N': out.Print("Ignoring addition of file ",f)
			}
		}
		cook.Undo(unraw)
	}
	if *shortlog == "COMMITNAME" {
		cook.SetCooked()
		out.Print("What is the patch name? ")
		inp,e := bufio.NewReaderSize(os.Stdin,1)
		error.FailOn(e)
		name,e := inp.ReadString('\n')
		error.FailOn(e)
		*shortlog = name
	}
	heads, _ := plumbing.ShowRef("--heads")
	hs := make([]git.Commitish,0,len(heads))
	for _,h := range heads {
		hs = hashes.Append(hs, h)
	}
	c := plumbing.CommitTree(plumbing.WriteTree(), hs, *shortlog)
	ctested, e := test.Commit(c)
	error.FailOn(e)
	plumbing.UpdateRef("HEAD", ctested)
}
