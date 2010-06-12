package main;

import (
	"os"
	"bufio"
	"goopt"
	git "./git/git"
	"./git/plumbing"
	"./util/out"
	"./util/debug"
	"./util/error"
	"./util/help"
	"./iolaus/prompt"
	"./iolaus/test"
	"./iolaus/core"
	hashes "./gotgo/slice(git.Commitish)"
)

var shortlog = goopt.String([]string{"-m","--patch"}, "COMMITNAME",
	"name of commit")

var description = func() string {
	return `
Record is used to name a set of changes and record the patch to the
repository.
`}

func main() {
	goopt.Vars["Verb"] = "Record"
	goopt.Vars["verb"] = "record"
	defer error.Exit(nil) // Must call exit so that cleanup will work!
	help.Init("record changes.", description, core.ModifiedFiles)
	git.AmInRepo("Must be in a repository to call record!")
	//plumbing.ReadTree(git.Ref("HEAD"))

	e := plumbing.ReadTree(git.Ref("HEAD"), "--index-output=.git/index.recording")
	if e != nil {
		error.Print("It looks like your repository is headless...")
	}
	plumbing.RefreshIndex()	// make sure the stat info is up-to-date...
	// Check which files are touched before using the new index...
	modfiles,e := core.DiffFiles([]string{})
	error.FailOn(e)
	// It's pretty hokey to use os.Setenv here rather than using exec to
	// set it directly, but it shouldn't be a problem as long as we
	// aren't calling git from multiple goroutines.
	e = os.Setenv("GIT_INDEX_FILE", ".git/index.recording")
	error.FailOn(e)

	prompt.Run(modfiles, func (f core.FileDiff) {
		plumbing.UpdateIndexCache(f.Info())
	})
	if *shortlog == "COMMITNAME" {
		out.Println("What is the patch name? ")
		inp,e := bufio.NewReaderSize(os.Stdin,1)
		error.FailOn(e)
		name,e := inp.ReadString('\n')
		error.FailOn(e)
		*shortlog = name
	}
	debug.Println("Reading heads...")
	heads, _ := plumbing.ShowRef("--heads")
	hs := make([]git.Commitish,0,len(heads))
	for _,h := range heads {
		hs = hashes.Append(hs, h)
	}
	debug.Println("Creating commit...")
	tree,e := plumbing.WriteTree()
	error.FailOn(e)
	c := plumbing.CommitTree(tree, hs, *shortlog)
	debug.Println("Testing commit...")
	ctested, e := test.Commit(c)
	error.FailOn(e)
	debug.Println("Updating HEAD...")
	plumbing.UpdateRef("HEAD", ctested)
	// Now let's update the true index by just copying over the scratch...
	e = os.Rename(".git/index.recording", ".git/index")
	error.FailOn(e) // FIXME: we should do better than this...
}
