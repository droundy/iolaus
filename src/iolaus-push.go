package main;

import (
	"github.com/droundy/goopt.git"
	git "./git/git"
	"./git/plumbing"
	"./util/out"
	"./util/exit"
	"./util/error"
	"./util/help"
	hashes "./gotgo/slice(git.Commitish)"
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

	remotes,e := plumbing.LsRemote("--heads", "origin")
	error.FailOn(e)
	locals,e := plumbing.ShowRef("--heads")
	error.FailOn(e)
	// Fetch the remotes so that they will be present when we need them later.
	origin := plumbing.RemoteUrl("origin")
	plumbing.FetchPack(origin, "--all", "-q")
	// Stick hashes into nice arrays...
	localrefs := make([]git.Commitish, 0, len(locals))
	for _,h := range locals {
		localrefs = hashes.Append(localrefs, h)
	}
	remoterefs := make([]git.Commitish, 0, len(remotes))
	for _,h := range remotes {
		remoterefs = hashes.Append(remoterefs, h)
	}
	topush, e := plumbing.RevListDifference(localrefs, remoterefs)
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
	topull, e := plumbing.RevListDifference(remoterefs, localrefs)
	error.FailOn(e)
	if len(topull) > 0 {
		for _,tp := range topull {
			cc, e := plumbing.Commit(tp)
			error.FailOn(e)
			out.Println("Could pull:\n", cc)
		}
		out.Println("I haven't finished with push yet.")
		exit.Exit(0)
	} else {
		out.Println("This is a fast-forward push!")
		if *all {
			error.FailOn(plumbing.SendPack(origin, locals))
		} else {
			out.Println("I haven't yet implemented interactive pushes.")
		}
	}
	exit.Exit(0)
}
