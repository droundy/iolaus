package plumbing

import (
	"./git"
	"os"
	"strings"
	"patch"
	"fmt"
)

func Init() {
	git.Run("init",nil)
}

type Patch patch.Set

func DiffFiles(paths []string) Patch {
	args := make([]string,len(paths)+2)
	args[0] = "-p"
	args[1] = "--"
	for i,p := range paths {
		args[i+2] = p
	}
	o, _ := git.Read("diff-files", args)
	p, _ := patch.Parse(strings.Bytes(o));
	return Patch(*p)
}

func (s Patch) String() (out string) {
	out = s.Header
	for _, f := range s.File {
		if len(out) > 0 { out += "\n" }
		out += string(f.Verb) + " " + f.Src
		if f.OldMode != 0 || f.NewMode != 0 {
			out += "\n" + fmt.Sprint(f.OldMode) + " -> " + fmt.Sprint(f.NewMode)
		}
		switch d := f.Diff.(type) {
		default:
			out += fmt.Sprintf("\nunexpected type %T", f.Diff)  // %T prints type
		case patch.TextDiff:
			for _, chunk := range d {
				out += "\n" + fmt.Sprintln(chunk.Line)
				older := strings.SplitAfter(string(chunk.Old), "\n",0)
				newer := strings.SplitAfter(string(chunk.New), "\n",0)
				for i, l := range older {
					if len(l)>0 {
						if i>=len(newer) || newer[i] != l {
							out += "-" + l
						} else {
							out += " " + l
						}
					}
				}
				for i, l := range newer {
					if len(l)>0 {
						if i>=len(older) || older[i] != l {
							out += "+" + l
						}
					}
				}
			}
		}
	}
	return
}

func LsFiles() (fs []string, e os.Error) {
	var o string
	o, e = git.Read("ls-files",
		              []string{"--exclude-standard", "-z", "--others", "--cached"})
	fs = strings.Split(o, "\000", 0)
	return
}
