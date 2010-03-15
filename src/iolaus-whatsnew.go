package main;

import (
	"strings"
	"io/ioutil"
	"./git/git"
	"./git/color"
	"./git/plumbing"
	"./util/out"
	"./util/help"
	"./util/exit"
	"./util/error"
	"./util/patience"
	"./iolaus/core"
)

var description = func() string {
	return `
whatsnew gives you a view of what changes you've made in your working
copy that haven't yet been recorded.
`}

func main() {
	help.Init("see unrecorded changes.", description, plumbing.LsFiles)
	git.AmInRepo("Must be in a repository to call whatsnew!")
	//plumbing.ReadTree(git.Ref("HEAD"))
	ds,e := core.DiffFiles([]string{})
	error.FailOn(e)
	// The following should all be factored apart next...
	for _,d := range ds {
		switch d.Change {
		case plumbing.Added:
			out.Print(color.String("Added "+d.Name, color.Meta))
		case plumbing.Deleted:
			out.Print(color.String("Deleted "+d.Name, color.Meta))
		case plumbing.Modified:
			if d.OldMode != d.NewMode { out.Println(d) }
			oldf, e := plumbing.Blob(d.OldHash)
			error.FailOn(e)
			newf, e := ioutil.ReadFile(d.Name)
			older := strings.SplitAfter(string(oldf),"\n",0)
			newer := strings.SplitAfter(string(newf),"\n",0)
			error.FailOn(e)
			mychunks := patience.Diff(older, newer)

			chunkLine := 1
			lastline := chunkLine
			if len(mychunks) == 0 { continue }
			if mychunks[0].Line-4 > lastline {
				lastline = mychunks[0].Line - 4
			}
			out.Printf(color.String("¤¤¤ %s %d ¤¤¤\n", color.Meta),d.Name,lastline)
			for _,ch := range mychunks {
				if ch.Line > lastline + 6 {
					for i:=lastline+1; i<lastline+4; i++ {
						out.Print(" ",newer[i-chunkLine])
					}
					out.Printf(color.String("¤¤¤ %s %d ¤¤¤\n", color.Meta),d.Name,ch.Line-3)
					for i:=ch.Line-3; i<ch.Line; i++ {
						out.Print(" ",newer[i-chunkLine])
					}
				} else {
					for i:=lastline; i<ch.Line; i++ {
						out.Print(" ",newer[i-chunkLine])
					}
				}
				out.Print(ch)
				lastline = ch.Line + len(ch.New)
			}
			for i:=lastline-chunkLine; i<len(newer)-1 && i < lastline-chunkLine+3;i++ {
				out.Print(" ",newer[i])
			}
		default:
			out.Println(d)
		}
	}
	exit.Exit(0)
}
