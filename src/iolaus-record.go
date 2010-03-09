package main;

import (
	"os"
	"bufio"
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
var all = goopt.Flag([]string{"-a","--all"}, []string{"--interactive"},
	"record all patches", "prompt for patches interactively")

func main() {
	git.AmInRepo("Must be in a repository to call record!")
	help.Init("record changes.", plumbing.LsFiles)
	//plumbing.ReadTree(plumbing.Ref("HEAD"))

	if *all {
		for _,f := range plumbing.DiffFilesModified([]string{}) {
			out.Print("Considering changes to",f)
			plumbing.UpdateIndex(f)
		}
		for _,newf := range plumbing.LsOthers() {
			out.Print("Considering adding",newf)
			plumbing.UpdateIndex(newf)
		}
	} else {
		defer cook.Undo(cook.SetRaw())
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
	}
	if *shortlog == "COMMITNAME" {
		defer cook.Undo(cook.SetCooked())
		out.Print("What is the patch name? ")
		inp,e := bufio.NewReaderSize(os.Stdin,1)
		error.FailOn(e)
		name,e := inp.ReadString('\n')
		error.FailOn(e)
		*shortlog = name
	}
	_, heads, _ := plumbing.ShowRef("--heads")
	c := plumbing.CommitTree(plumbing.WriteTree(), heads, *shortlog)
	plumbing.UpdateRef("HEAD", c)
}
