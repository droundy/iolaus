package main;

import (
	"os"
	"fmt"
	"github.com/droundy/goopt"
	git "./git/git"
	"./git/plumbing"
	"./util/out"
	"./util/debug"
	"./util/exit"
	"./util/error"
	"./util/help"
	"./iolaus/test"
	"./iolaus/core"
	"./iolaus/promptcommit"
)

var description = func() string {
	return `
Pull is used to bring changes made in another repository into the current
repository (that is, either the one in the current directory, or the one
specified with the --repodir option). Pull allows you to bring over all or
some of the patches that are in that repository but not in this one. Pull
accepts arguments, which are URLs from which to pull, and when called
without an argument, pull will pull from origin.
`}

func main() {
	goopt.Vars["Verb"] = "Pull"
	goopt.Vars["verb"] = "pull"
	help.Init("pull changes from origin.", description, plumbing.LsFiles)
	git.AmInRepo("Must be in a repository to call pull!")

	origin := plumbing.RemoteUrl("origin")
	if len(goopt.Args) > 0 {
		origin = plumbing.RemoteUrl(goopt.Args[0])
		out.Println("Pulling from ", origin)
	}
	remote,e := plumbing.RemoteHead(origin)
	error.FailOn(e)
	// ignore error code on show-ref, since it returns an error when
	// there are no heads.
	local,_ := plumbing.LocalHead()
	// Fetch the remotes so that they will be present when we need them later.
	plumbing.FetchPack(origin, "--all", "-q")
	debug.Println("Looking for stuff to push...")
	topush, e := plumbing.RevListDifference([]git.Commitish{local}, []git.Commitish{remote})
	error.FailOn(e)

	if len(topush) > 0 {
		// We need to do a real merge!
		for _,tp := range topush {
			cc, e := plumbing.Commit(tp)
			error.FailOn(e)
			out.Println("Could push:\n", cc)
		}
		// barf on local changes...
		if !local.IsEmpty() {
			plumbing.RefreshIndex()
			p,e := plumbing.DiffFiles([]string{})
			error.FailOn(e)
			if len(p) > 0 {
				error.FailOn(os.NewError(
					fmt.Sprintf("I can't handle local changes yet! %v",p)))
			}
		}
		// It's pretty hokey to use os.Setenv here rather than using exec to
		// set it directly, but it shouldn't be a problem as long as we
		// aren't calling git from multiple goroutines.
		error.FailOn(plumbing.ReadTree(local, "--index-output=.git/index.pulling"))
		error.FailOn(os.Setenv("GIT_INDEX_FILE", ".git/index.pulling"))
		t, e := core.Merge(local, remote)
		error.FailOn(e)
		c := plumbing.CommitTree(t, []git.Commitish{local,remote}, "Merge")
		debug.Println("Testing merge...")
		ctested, e := test.Commit(c)
		error.FailOn(e)
		plumbing.UpdateRef("HEAD", ctested)
		plumbing.CheckoutIndex("--force", "--all", "-u")
		// Now let's update the true index by just copying over the scratch...
		error.FailOn(os.Rename(".git/index.pulling", ".git/index"))
		exit.Exit(0)
	} else {
		out.Println("This is a fast-forward pull!")

		debug.Println("Prompting for commits...")
		hnew := promptcommit.Select(local, remote)
		if local.IsEmpty() {
			debug.Println("We have no local commits...")
			e = plumbing.ReadTree(hnew)
		} else {
			debug.Println("Merging with local index...")
			e = plumbing.ReadTree2(local, hnew)
		}
		error.FailOn(e)
		plumbing.UpdateRef("HEAD", hnew)
		plumbing.CheckoutIndex("--force", "--all")
		exit.Exit(0)
	}
	exit.Exit(0)
}
