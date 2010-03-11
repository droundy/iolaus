package main;

import (
	"goopt"
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

func main() {
	git.AmInRepo("Must be in a repository to call push!")
	help.Init("push changes to origin.", plumbing.LsFiles)

	remotes,e := plumbing.LsRemote("origin", "--heads")
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
		out.Print("Could push:\n", cc)
	}
	if len(topush) == 0 {
		out.Print("No commits to push!")
		exit.Exit(0)
	}
	if *dryRun { exit.Exit(0) }
	topull, e := plumbing.RevListDifference(remoterefs, localrefs)
	error.FailOn(e)
	if len(topull) > 0 {
		for _,tp := range topull {
			cc, e := plumbing.Commit(tp)
			error.FailOn(e)
			out.Print("Could pull:\n", cc)
		}
		out.Print("I haven't finished with push yet.")
		exit.Exit(0)
	} else {
		out.Print("This is a fast-forward push!")
		if *all {
			plumbing.SendPack(origin, locals)
		} else {
			out.Print("I haven't yet implemented interactive pushes.")
		}
	}
	exit.Exit(0)
}
