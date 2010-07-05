package promptcommit

import (
	"github.com/droundy/goopt"
	"../git/git"
	"../git/plumbing"
	"../util/out"
	"../util/error"
	"../util/exit"
	"../util/debug"
	box "./gotgo/box(git.CommitHash,git.Commitish)"
)

// To make --all the default, set *prompt.All to true.
var All = goopt.Flag([]string{"-a","--all"}, []string{"--interactive"},
	"verb all commits", "prompt for commits interactively")

var Dryrun = goopt.Flag([]string{"--dry-run"}, []string{},
	"just show commits that would be verbed", "xxxs")

func Select(since, upto git.Commitish) (outh git.CommitHash) {
	hs,e := plumbing.RevListDifference([]git.Commitish{upto}, []git.Commitish{since})
	error.FailOn(e)
	if len(hs) == 0 {
		out.Println(goopt.Expand("No commits to verb!"))
		exit.Exit(0)
	}
	if *Dryrun {
		for _,h := range hs {
			cc, e := plumbing.Commit(h)
			error.FailOn(e)
			out.Println(goopt.Expand("Could verb:\n"), cc)
		}
		exit.Exit(0)
	}
	if *All {
		debug.Println(goopt.Expand("Verbing --all commits..."))
		h,e := plumbing.GetCommitHash(upto)
		error.FailOn(e)
		return h
	} else {
	  for len(hs) > 0 {
			// Just keep asking until we get a reasonable answer...
			cc,e := plumbing.Commit(hs[0])
			error.FailOn(e)
			out.Println(cc)
			c,e := out.PromptForChar(goopt.Expand("Verb this commit? "))
			error.FailOn(e)
			switch c {
			case 'q','Q': error.Exit(e)
			case 'v','V':
				c,e := plumbing.Commit(hs[0])
				error.FailOn(e)
				out.Println(c)
			case 'y','Y':
				debug.Println(goopt.Expand("Verbing hash "),hs[0])
				outh = hs[0]
				hs,e = plumbing.RevListDifference(box.Box(hs[1:]), []git.Commitish{hs[0]})
				error.FailOn(e)
				continue
			case 'n','N': debug.Println("Ignoring commit ",hs[0])
				out.Println("BUGGY CODE HERE!!!")
				// this is wrong...
				hs,e = plumbing.RevListDifference(box.Box(hs[1:]), []git.Commitish{hs[0]})
				error.FailOn(e)
				continue
		}
		}
	}
	return
}
