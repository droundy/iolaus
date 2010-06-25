package main;

import (
	"github.com/droundy/goopt.git"
	"./git/git"
	"./git/plumbing"
	"./util/help"
	"./util/out"
	"./util/error"
	"./util/exit"
)

var description = func() string {
	return `
show information about commits.
`}

func main() {
	goopt.Vars["Verb"] = "Display"
	goopt.Vars["verb"] = "display"
	help.Init("see commits.", description, plumbing.LsFiles)
	git.AmInRepo("Must be in a repository to call changes!")
	heads, _ := plumbing.ShowRef("--heads")
	hs := make([]git.CommitHash,len(heads))
	i := 0
	for _,v := range heads {
		hs[i] = v
	}
	showChanges(hs, make(map[string]bool))
	exit.Exit(0)
}

func showChanges(hs []git.CommitHash, done map[string]bool) {
	if len(hs) == 0 { return }
	newhs := make([]git.CommitHash, 0)
	for _,h := range hs {
		if _,amdone := done[h.String()]; !amdone {
			c,e := plumbing.Commit(h)
			error.FailOn(e)
			done[h.String()] = true
			out.Println(c)
			newhs = csCat(newhs, c.Parents)
		}
	}
	showChanges(newhs, done)
}
