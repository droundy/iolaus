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
	"pull all patches", "prompt for patches interactively")
var dryRun = goopt.Flag([]string{"--dry-run"}, []string{},
	"don't actually pull, just show what we would pull",
	"actually do pull")

func main() {
	git.AmInRepo("Must be in a repository to call pull!")
	help.Init("pull changes from origin.", plumbing.LsFiles)

	origin := plumbing.RemoteUrl("origin")
	if len(goopt.Args) > 1 {
		origin = plumbing.RemoteUrl(goopt.Args[1])
		out.Print("Pulling from ", origin)
	}
	remotes,e := plumbing.LsRemote(origin, "--heads")
	error.FailOn(e)
	// ignore error code on show-ref, since it returns an error when
	// there are no heads.
	locals,_ := plumbing.ShowRef("--heads")
	// Fetch the remotes so that they will be present when we need them later.
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
	topull, e := plumbing.RevListDifference(remoterefs, localrefs)
	error.FailOn(e)
	for _,tp := range topull {
		cc, e := plumbing.Commit(tp)
		error.FailOn(e)
		out.Print("Could pull:\n", cc)
	}
	if len(topull) == 0 {
		out.Print("No commits to pull!")
		exit.Exit(0)
	}
	if *dryRun { exit.Exit(0) }
	if len(topush) > 0 {
		for _,tp := range topush {
			cc, e := plumbing.Commit(tp)
			error.FailOn(e)
			out.Print("Could push:\n", cc)
		}
		out.Print("I haven't finished with pull yet.")
		exit.Exit(0)
	} else {
		out.Print("This is a fast-forward pull!")
		if *all {
			//plumbing.SendPack(origin, locals)
			if len(remotes) == 1 {
				for _,h := range remotes {
					plumbing.UpdateRef("HEAD", h)
					plumbing.ReadTree(h)
					plumbing.CheckoutIndex("--all")
					exit.Exit(0)
				}
			} else {
				out.Print("I haven't yet implemented merging.")
			}
		} else {
			out.Print("I haven't yet implemented interactive pulls.")
		}
	}
	exit.Exit(0)
}
