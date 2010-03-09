package plumbing

import (
	"./git"
	"../util/patience"
	"../util/debug"
	"../util/error"
	"os"
	"strings"
	"patch"
	"fmt"
	stringslice "./gotgo/slice(string)"
)

func Init() {
	git.Run("init")
}

type Patch patch.Set

type Hash [40]byte
func (h Hash) String() string { return string(h[0:40]) }
type Tree string
func (r Tree) String() string { return string(r) }
type Ref string
func (r Ref) String() string { return string(r) }

type Treeish interface {
	String() string
}

type Commitish interface {
	String() string
}

func UpdateIndex(f string) os.Error {
	return git.Run("update-index", "--add", "--remove", "--", f)
}

func UpdateRef(ref string, val Commitish) os.Error {
	return git.Run("update-ref", ref, val.String())
}

func WriteTree() Treeish {
	o,_ := git.Read("write-tree")
	return Tree(o[0:40])
}

func CommitTree(tree Treeish, parents []Commitish, log string) Commitish {
	args := []string{tree.String()}
	for _,p := range parents {
		args = stringslice.Append(args, "-p")
		args = stringslice.Append(args, p.String())
	}
	o,e := git.WriteReadS("commit-tree", log, args)
	if e != nil { panic("bad output in commit-tree") }
	return Ref(o[0:40])
}

func ReadTree(ref Treeish) {
	git.Run("read-tree", ref.String())
}

func DiffFilesModified(paths []string) []string {
	args := stringslice.Cat([]string{"--name-only", "-z", "--"}, paths)
	o,_ := git.ReadS("diff-files", args)
	return splitOnNulls(o)
}

func DiffFiles(paths []string) Patch {
	args := stringslice.Cat([]string{"-p", "--"}, paths)
	o,e := git.ReadS("diff-files", args)
	error.FailOn(e)
	p, e := patch.Parse([]byte(o));
	error.FailOn(e)
	return Patch(*p)
}

func (s Patch) String() (out string) {
	out = s.Header
	for _, f := range s.File {
		if len(out) > 0 { out += "\n" }
		// out += string(f.Verb) + " " + f.Src
		if f.OldMode != 0 && f.NewMode != 0 {
			out += "\n" + fmt.Sprint(f.OldMode) + " -> " + fmt.Sprint(f.NewMode)
		}
		if f.OldMode != 0 && f.NewMode == 0 { // The file has been removed!
			out += "\nRemoved "+ f.Src
			continue
		}
		switch d := f.Diff.(type) {
		default:
			if f.Diff == patch.NoDiff {
				// There is nothing here to see...
			} else {
				out += fmt.Sprintf("\nunexpected type %T", f.Diff)  // %T prints type
			}
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
	return genLsFilesE("--exclude-standard","-z","--others","--cached")
}

func LsOthers() []string {
	o, e := genLsFilesE("--exclude-standard","-z","--others")
	if e != nil {
		debug.Print("yeckgys")
	}
	return o
}

func genLsFilesE(args ...string) ([]string, os.Error) {
	o, e := git.Read("ls-files", args)
	return splitOnNulls(o), e
}

func splitOnNulls(s string) []string {
	xs := strings.Split(s, "\000", 0)
	if len(xs[len(xs)-1]) == 0 {
		return xs[0:len(xs)-1]
	}
	return xs
}
