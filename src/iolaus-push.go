package main;

import (
	"goopt"
	"./git/git"
	"./git/plumbing"
	"./util/out"
	"./util/exit"
	"./util/error"
	"./util/help"
)

var all = goopt.Flag([]string{"-a","--all"}, []string{"--interactive"},
	"push all patches", "prompt for patches interactively")

func main() {
	git.AmInRepo("Must be in a repository to call push!")
	help.Init("push changes to origin.", plumbing.LsFiles)

	remotes, _, e := plumbing.LsRemote("origin", "--heads")
	// Fetch the remotes asynchronously in the hopes that they will be
	// present when we need them later.
	go plumbing.FetchPack(plumbing.RemoteUrl("origin"), "--all", "-q")
	error.FailOn(e)
	out.Print("Remotes are ", remotes)
	out.Print("I haven't finished with push yet.")
	for _,r := range remotes {
		cc, e := plumbing.Commit(r)
		error.FailOn(e)
		out.Print("Remote has:\n", cc)
	}
	exit.Exit(0)
}
