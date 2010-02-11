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

func Diff(paths []string) Patch {
	o, _ := git.Read("diff", []string{})
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
				numstart := 0
				for i:=0; i<5; i++ {
					if i < len(older)-1 && i < len(newer)-1 &&
						older[i] == newer[i] {
						numstart = i+1
					} else {
						break
					}
				}
				numend := 0
				for i:=1; i<5; i++ {
					if i < len(older) && i < len(newer) &&
						older[len(older)-i] == newer[len(newer)-i] {
						numend = i
					} else {
						break
					}
				}
				out += "hello world "+ fmt.Sprint(numstart)+"\n"
				for _, l := range older[0:numstart] {
					out += " " + l
				}
				for _, l := range older[numstart:len(older)-numend] {
					out += "-" + l
				}
				for _, l := range newer[numstart:len(newer)-numend] {
					out += "+" + l
				}
				for _, l := range older[len(older)-numend:len(older)] {
					out += " " + l
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
