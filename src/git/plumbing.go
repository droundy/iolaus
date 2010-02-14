package plumbing

import (
	"./git"
	"../util/patience"
	"../util/debug"
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
		// out += string(f.Verb) + " " + f.Src
		if f.OldMode != 0 || f.NewMode != 0 {
			out += "\n" + fmt.Sprint(f.OldMode) + " -> " + fmt.Sprint(f.NewMode)
		}
		switch d := f.Diff.(type) {
		default:
			out += fmt.Sprintf("\nunexpected type %T", f.Diff)  // %T prints type
		case patch.TextDiff:
			for _, chunk := range d {
				older := strings.SplitAfter(string(chunk.Old), "\n",0)
				newer := strings.SplitAfter(string(chunk.New), "\n",0)
				lastline:=chunk.Line
				mychunks := patience.DiffFromLine(chunk.Line, older, newer)
				fmt.Println("¤¤¤¤¤¤", f.Src,chunk.Line,"¤¤¤¤¤¤")
				for _,ch := range mychunks {
					if ch.Line > lastline + 6 {
						for i:=lastline+1; i<lastline+4; i++ {
							fmt.Print(" ",newer[i-chunk.Line])
						}
						fmt.Println("¤¤¤¤¤¤", f.Src,ch.Line-3,"¤¤¤¤¤¤")
						for i:=ch.Line-3; i<ch.Line; i++ {
							fmt.Print(" ",newer[i-chunk.Line])
						}
					} else {
						for i:=lastline; i<ch.Line; i++ {
							fmt.Print(" ",newer[i-chunk.Line])
						}
					}
					fmt.Print(ch)
					lastline = ch.Line + len(ch.New)
				}
				for i:=lastline-chunk.Line; i<len(newer);i++ {
					fmt.Print(" ",newer[i])
				}
			}
		}
	}
	return
}

func LsFiles() []string {
	o, _ := LsFilesE()
	return o
}

func LsFilesE() (fs []string, e os.Error) {
	return genLsFilesE([]string{"--exclude-standard","-z","--others","--cached"})
}

func LsOthers() []string {
	o, e := genLsFilesE([]string{"--exclude-standard","-z","--others"})
	if e != nil {
		debug.Print("yeckgys")
	}
	return o
}

func genLsFilesE(args []string) ([]string, os.Error) {
	o, e := git.Read("ls-files", args)
	fs := strings.Split(o, "\000", 0)
	debug.Print("ls-files gives",o)
	if len(fs[len(fs)-1]) == 0 {
		return fs[0:len(fs)-1], e
	}
	return fs, e
}
