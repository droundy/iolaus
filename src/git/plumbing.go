package plumbing

import (
	git "./git"
	"../util/patience"
	"../util/debug"
	"../util/error"
	"os"
	"once"
	"strings"
	"patch"
	"fmt"
	stringslice "./gotgo/slice(string)"
	pars "./gotgo/slice(git.CommitHash)"
)

func Init() {
	git.Run("init")
}

type Patch patch.Set

func UpdateIndex(f string) os.Error {
	return git.Run("update-index", "--add", "--remove", "--", f)
}

func CheckoutIndex(args ...string) os.Error {
	return git.RunS("checkout-index", args)
}

func UpdateRef(ref string, val git.Commitish) os.Error {
	return git.Run("update-ref", ref, val.String())
}

func FetchPack(remote string, args ...string) {
	args = stringslice.Append(args, remote)
	git.RunSilentlyS("fetch-pack", args)
}

func LsRemote(remote string, args ...string) (hs map[git.Ref]git.CommitHash, e os.Error) {
	args = stringslice.Append(args, remote)
	o, e := git.ReadS("ls-remote", args)
	if e != nil { return }
	return splitRefs(o), e
}

func ShowRef(args ...string) (hs map[git.Ref]git.CommitHash, e os.Error) {
	o, e := git.Read("show-ref", args)
	if e != nil { return }
	return splitRefs(o), e
}

func splitRefs(s string) (hs map[git.Ref]git.CommitHash) {
	hs = make(map[git.Ref]git.CommitHash)
	xs := strings.Split(s, "\n", 0)
	for _,x := range xs {
		if len(x) > 42 {
			hs[git.Ref(string(x[41:]))] = git.CommitHash(mkHash(x))
		}
	}
	return
}

func WriteTree() git.Treeish {
	o,_ := git.Read("write-tree")
	t := git.TreeHash{}
	for j := range o[0:40] {
		t[j] = o[j] // shouldn't there be a nicer way to do this?
	}
	return t
}

func CommitTree(tree git.Treeish, parents []git.Commitish, log string) git.Commitish {
	args := []string{tree.String()}
	for _,p := range parents {
		args = stringslice.Append(args, "-p")
		args = stringslice.Append(args, p.String())
	}
	o,e := git.WriteReadS("commit-tree", log, args)
	if e != nil { panic("bad output in commit-tree: "+e.String()) }
	return git.Ref(o[0:40]) // should be a git.Commitish
}

func ReadTree(ref git.Treeish) {
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

var configList map[string]string

func ListConfig() map[string]string {
	once.Do(func() {
		o,_ := git.Read("config", "--null", "--list")
		configList = splitOnNullsAndLines(o)
	})
	return configList
}

func splitOnNullsAndLines(s string) (out map[string]string) {
	xs := strings.Split(s, "\000", 0)
	out = make(map[string]string)
	for _,x := range xs {
		kv := strings.Split(x, "\n", 0)
		if len(kv) == 2 {
			out[kv[0]] = kv[1]
		}
	}
	return
}

func RemoteUrl(r string) string {
	conf := ListConfig()
	rr,ok := conf["remote."+r+".url"]
	if ok { return rr }
	return r // not a real remote
}

type CommitEntry struct {
	Parents []git.CommitHash
  Tree git.TreeHash
  Author string
  Committer string
  Message string
}

func (ce CommitEntry) String() string {
	return "Author: " + ce.Author + "\n" + ce.Message
}

func rawCommit(c git.Commitish) (o string, e os.Error) {
	return git.Read("cat-file", "commit", c.String())
}

func Commit(c git.Commitish) (ce CommitEntry, e os.Error) {
	o,e := rawCommit(c)
	if e != nil { return }
	ls := strings.Split(o, "\n", 0)
	// ls will always have a length of at least one!
	for _,l := range ls {
		switch {
		case startsWith(l, "tree "):
			ce.Tree = git.TreeHash(mkHash(dropToSpace(l)))
		case startsWith(l, "parent "):
			ce.Parents = pars.Append(ce.Parents, git.CommitHash(mkHash(dropToSpace(l))))
		case startsWith(l, "author "):
			ce.Author = dropToSpace(l)
		case startsWith(l, "committer "):
			ce.Committer = dropToSpace(l)
		default:
			break
		}
	}
	xs := strings.Split(o, "\n\n", 2)
	if len(xs) != 2 { panic("bad foo in bar") }
	ce.Message = xs[1]
	return
}

func startsWith(s string, e string) bool {
	return len(s) >= len(e) && s[0:len(e)] == e
}

func dropToSpace(s string) string {
	for i,c := range s {
		if c == ' ' {
			return s[i+1:]
		}
	}
	return ""
}

func mkHash(s string) (h git.Hash) {
	if len(s) < 40 { panic("Hash to small in mkhash") }
	for i,v := range []byte(s)[0:40] { h[i] = v }
	return
}

func RevListDifference(newer, older []git.Commitish) ([]git.CommitHash, os.Error) {
	args := []string{}
	for _,n := range newer {
		args = stringslice.Append(args, n.String())
	}
	for _,o := range older {
		args = stringslice.Append(args, "^"+o.String())
	}
	return RevListS(args)
}

func RevListS(args []string) (commits []git.CommitHash, e os.Error) {
	o,e := git.ReadS("rev-list", args)
	if e != nil { return }
	ls := strings.Split(o,"\n",0)
	for _,l := range ls {
		if len(l) == 40 {
			commits = pars.Append(commits, git.CommitHash(mkHash(l)))
		}
	}
	return
}

func SendPack(repo0 string, updates map[git.Ref]git.CommitHash) os.Error {
	repo := RemoteUrl(repo0)
	args := []string{repo}
	for name,hash := range updates {
		args = stringslice.Append(args, hash.String()+":"+name.String())
	}
	return git.RunS("send-pack", args)
}
