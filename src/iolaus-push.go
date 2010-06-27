package main;

import (
	"os"
	"github.com/droundy/goopt"
	git "./git/git"
	"./git/plumbing"
	"./util/out"
	"./util/debug"
	"./util/exit"
	"./util/error"
	"./util/help"
	"./iolaus/core"
	"./iolaus/test"
)

var all = goopt.Flag([]string{"-a","--all"}, []string{"--interactive"},
	"push all patches", "prompt for patches interactively")
var dryRun = goopt.Flag([]string{"--dry-run"}, []string{},
	"don't actually push, just show what we would push",
	"actually do push")

var description = func() string {
	return `
Push is the opposite of pull.  Push allows you to copy changes from the
current repository into another repository.
`}

func main() {
	help.Init("push changes to origin.", description, plumbing.LsFiles)
	git.AmInRepo("Must be in a repository to call push!")

	// Fetch the remotes so that they will be present when we need them later.
	origin := plumbing.RemoteUrl("origin")
	plumbing.FetchPack(origin, "--all", "-q")
	// We use the remote "master", since for some reason we don't seem
	// to be able to push to a remote "head".  :(
	remote,e := plumbing.RemoteMaster(origin)
	out.Println("remote is", remote)
	error.FailOn(e)
	// look up the local head
	local,e := plumbing.LocalHead()
	error.FailOn(e)
	topush, e := plumbing.RevListDifference([]git.Commitish{local}, []git.Commitish{remote})
	error.FailOn(e)
	for _,tp := range topush {
		cc, e := plumbing.Commit(tp)
		error.FailOn(e)
		out.Println("Could push:\n", cc)
	}
	if len(topush) == 0 {
		out.Println("No commits to push!")
		exit.Exit(0)
	}
	if *dryRun { exit.Exit(0) }
	topull, e := plumbing.RevListDifference([]git.Commitish{remote}, []git.Commitish{local})
	error.FailOn(e)
	if len(topull) > 0 {
		for _,tp := range topull {
			cc, e := plumbing.Commit(tp)
			error.FailOn(e)
			out.Println("Could pull:\n", cc)
		}
		// It's pretty hokey to use os.Setenv here rather than using exec to
		// set it directly, but it shouldn't be a problem as long as we
		// aren't calling git from multiple goroutines.
		error.FailOn(plumbing.ReadTree(local, "--index-output=.git/index.pushing"))
		error.FailOn(os.Setenv("GIT_INDEX_FILE", ".git/index.pushing"))
		t, e := core.Merge(local, remote)
		error.FailOn(e)
		c := plumbing.CommitTree(t, []git.Commitish{local,remote}, "Merge")
		debug.Println("Testing merge...")
		ctested, e := test.Commit(c)
		error.FailOn(e)
		// No need to check for success on removing the temporary index...
		os.Remove(".git/index.pushing")
		error.FailOn(plumbing.SendPack(origin,
			map[git.Ref]git.CommitHash{"refs/heads/master": ctested}))
		exit.Exit(0)
	} else {
		out.Println("This is a fast-forward push!")
		if *all {
			error.FailOn(plumbing.SendPack(origin,
				map[git.Ref]git.CommitHash{"refs/heads/master": local}))
		} else {
			out.Println("I haven't yet implemented interactive pushes.")
		}
	}
	exit.Exit(0)
}
