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
	"./iolaus/core"
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
	help.Init("record changes.", description, core.ModifiedFiles)
	git.AmInRepo("Must be in a repository to call record!")
	//plumbing.ReadTree(git.Ref("HEAD"))

	e := plumbing.ReadTree(git.Ref("HEAD"), "--index-output=.git/index.recording")
	if e != nil {
		error.Print("It looks like your repository is headless...")
	}
	// Check which files are touched before using the new index...
	modfiles := core.ModifiedFiles()
	// It's pretty hokey to use os.Setenv here rather than using exec to
	// set it directly, but it shouldn't be a problem as long as we
	// aren't calling git from multiple goroutines.
	e = os.Setenv("GIT_INDEX_FILE", ".git/index.recording")
	error.FailOn(e)

	if *all {
		for _,f := range modfiles {
			out.Println("Considering changes to ",f)
			plumbing.UpdateIndex(f)
		}
	} else {
		unraw := cook.SetRaw()
		defer cook.Undo(unraw)
		for _,f := range modfiles {
			c,e := out.PromptForChar("Record changes to %s? ", f)
			switch c {
			case 'q','Q': error.Exit(e)
			case 'y','Y':
				out.Println("Dealing with file ",f)
				plumbing.UpdateIndex(f)
			case 'n','N': out.Println("Ignoring changes to file ",f)
			}
		}
		cook.Undo(unraw)
	}
	if *shortlog == "COMMITNAME" {
		cook.SetCooked()
		out.Println("What is the patch name? ")
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
	// Now let's update the true index by just copying over the scratch...
	e = os.Rename(".git/index.recording", ".git/index")
	error.FailOn(e) // FIXME: we should do better than this...
}
