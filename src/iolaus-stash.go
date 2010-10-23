package main;

import (
	"os"
	"github.com/droundy/goopt"
	git "./git/git"
	"./git/porcelain"
	"./git/plumbing"
	"./util/error"
	"./util/help"
	"./iolaus/prompt"
	"./iolaus/core"
)

var stashname = goopt.String([]string{"-m"}, "STASHNAME",
	"name of stash")

var description = func() string {
	return `
Stash is used to undo a set of unrecorded changes, much like git stash.
`}

func main() {
	goopt.Vars["Verb"] = "Keep"
	goopt.Vars["verb"] = "keep"
	defer error.Exit(nil) // Must call exit so that cleanup will work!
	help.Init("record changes.", description, core.ModifiedFiles)
	git.AmInRepo("Must be in a repository to call record!")
	//plumbing.ReadTree(git.Ref("HEAD"))

	e := plumbing.ReadTree(git.Ref("HEAD"), "--index-output=.git/index.reverting")
	defer os.Remove(".git/index.reverting")
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
	e = os.Setenv("GIT_INDEX_FILE", ".git/index.reverting")
	error.FailOn(e)

	prompt.Run(modfiles, func (f core.FileDiff) {
		plumbing.UpdateIndexCache(f.Info())
	})
	stashargs := make([]string, 2, 3)
	stashargs[0] = "save"
	stashargs[1] = "--keep-index"
	if *stashname != "STASHNAME" {
		stashargs = stashargs[:3]
		stashargs[2] = *stashname
	}
	error.Exit(porcelain.Stash(stashargs...))
}
