package plumbing

import (
	git "./git"
	"./color"
	"../util/patience"
	"../util/debug"
	"../util/error"
	"os"
	"sync"
	"strings"
	"patch"
	"fmt"
	"strconv"
	stringslice "../util/slice(string)"
	pars "./gotgo/slice(git.CommitHash)"
)

func Init() {
	git.Run("init")
}

type Patch patch.Set

func UpdateIndex(f string) os.Error {
	return git.Run("update-index", "--add", "--remove", "--", f)
}

func RefreshIndex() os.Error {
	return git.Run("update-index", "--refresh")
}

func HashObject(c string) (h git.Hash, e os.Error) {
	o,e := git.WriteRead("hash-object", c, "-w", "--stdin")
	if e != nil { return }
	return mkHash(string(o)), e
}

func HashFile(f string) (h git.Hash, e os.Error) {
	o,e := git.Read("hash-object", "-w", f)
	if e != nil { return }
	return mkHash(string(o)), e
}

func HashObjectUpdateIndex(mode int, contents, f string) os.Error {
	h,e := HashObject(contents)
	if e != nil { return e }
	return UpdateIndexCache(mode, h, f)
}

func UpdateIndexCache(mode int, h git.Hash, f string) os.Error {
	if mode == 0 {
		return git.Run("update-index", "--force-remove", f)
	}
	return git.Run("update-index", "--add", "--cacheinfo",
		fmt.Sprintf("%o", mode), h.String(), f)
}

func CheckoutIndex(args ...string) os.Error {
	return git.Run("checkout-index", args...)
}

func UpdateRef(ref string, val git.Commitish) os.Error {
	return git.Run("update-ref", ref, val.String())
}

func FetchPack(remote string, args ...string) {
	args = stringslice.Append(args, remote)
	git.RunSilently("fetch-pack", args...)
}

func LsRemote(args ...string) (hs map[git.Ref]git.CommitHash, e os.Error) {
	o, e := git.ReadS("ls-remote", args)
	if e != nil { return }
	return splitRefs(o), e
}

func RemoteHead(remote string) (h git.CommitHash, e os.Error) {
	xs, e := LsRemote(remote)
	if e != nil { return }
	return xs["HEAD"], nil
}

func RemoteMaster(remote string) (h git.CommitHash, e os.Error) {
	xs, e := LsRemote(remote)
	if e != nil { return }
	return xs["refs/heads/master"], nil
}

func ShowRef(args ...string) (hs map[git.Ref]git.CommitHash, e os.Error) {
	o, e := git.Read("show-ref", args...)
	if e != nil { return }
	return splitRefs(o), e
}

func LocalHead() (h git.CommitHash, e os.Error) {
	return RevParse("HEAD")
}

func splitRefs(s string) (hs map[git.Ref]git.CommitHash) {
	hs = make(map[git.Ref]git.CommitHash)
	xs := strings.Split(s, "\n", -1)
	for _,x := range xs {
		if len(x) > 42 {
			hs[git.Ref(string(x[41:]))] = git.CommitHash(mkHash(x))
		}
	}
	return
}

func WriteTree() (h git.TreeHash, e os.Error) {
	o,e := git.Read("write-tree")
	if e != nil { return }
	return git.TreeHash(mkHash(o)), nil
}

func CommitTree(tree git.Treeish, parents []git.Commitish, log string) git.CommitHash {
	args := []string{tree.String()}
	for _,p := range parents {
		args = stringslice.Append(args, "-p")
		args = stringslice.Append(args, p.String())
	}
	o,e := git.WriteReadS("commit-tree", log, args)
	if e != nil { panic("bad output in commit-tree: "+e.String()) }
	return git.CommitHash(mkHash(o))
}

func ReadTree(ref git.Treeish, args ...string) os.Error {
	args = stringslice.Append(args, ref.String())
	return git.Run("read-tree", args...)
}

func ReadTree2(us, them git.Treeish, args ...string) os.Error {
	foo := args
	args = stringslice.Cat(foo, []string{"-m", us.String(), them.String()})
	return git.Run("read-tree", args...)
}

func ReadTree3(base, us, them git.Treeish, args ...string) os.Error {
	foo := args
	args = stringslice.Cat(foo,
		[]string{"-m", base.String(), us.String(), them.String()} )
	return git.Run("read-tree", args...)
}

func DiffFilesModified(paths []string) []string {
	args := stringslice.Cat([]string{"--name-only", "-z", "--"}, paths)
	o,_ := git.ReadS("diff-files", args)
	return splitOnNulls(o)
}

func DiffFilesP(paths []string) Patch {
	args := stringslice.Cat([]string{"-p", "--"}, paths)
	o,e := git.ReadS("diff-files", args)
	error.FailOn(e)
	p, e := patch.Parse([]byte(o));
	error.FailOn(e)
	return Patch(*p)
}

func DiffFiles(paths []string) (d []FileDiff, e os.Error) {
	args := stringslice.Cat([]string{"--"}, paths)
	o,e := git.ReadS("diff-files", args)
	if e != nil { return }
	ls := strings.Split(o,"\n",-1)
	d = make([]FileDiff, 0, len(ls))
	i := 0
	for _,l := range ls {
		if len(l) < 84 { continue }
		d = d[0:i+1]
		if l[0] != ':' { continue }
		// It would be faster to just use byte offsets, but I'm sure I'd
		// get them wrong, so for now I'll just use the slow, sloppy, lazy
		// approach.
		xxxx := strings.Split(l[1:], "\t",-1)
		if len(xxxx) < 2 {
			e = os.NewError("bad line: "+l)
			return
		}
		chunks := strings.Split(xxxx[0]," ",-1)
		if len(chunks) < 4 { continue }
		mode, e := strconv.Btoui64(chunks[0],8)
		if e != nil { return }
		d[i].OldMode = int(mode)
		mode, e = strconv.Btoui64(chunks[1],8)
		if e != nil { return }
		d[i].NewMode = int(mode)
		d[i].OldHash = mkHash(chunks[2])
		d[i].NewHash = mkHash(chunks[3])
		d[i].Change = Verb(chunks[4][0])
		d[i].Name = xxxx[1]
		switch d[i].Change {
		case Renamed, Copied:
			d[i].OldName = xxxx[2]
		}
		i++
	}
	return
}

type FileDiff struct {
	OldName, Name string
	OldMode, NewMode int
	OldHash, NewHash git.Hash
	Change Verb
}
func (d FileDiff) String() string {
	switch d.Change {
	case Copied, Renamed:
		return fmt.Sprintf(":%o %o %s %s %c\t%s\t%s",
			d.OldMode, d.NewMode, d.OldHash, d.NewHash, d.Change, d.Name, d.OldName)
	}
	return fmt.Sprintf(":%06o %06o %s %s %c\t%s",
		d.OldMode, d.NewMode, d.OldHash, d.NewHash, d.Change, d.Name)
}

type Verb int
const (
	Added Verb = 'A'
	Copied Verb = 'C'
	Deleted Verb = 'D'
	Modified Verb = 'M'
	Renamed Verb = 'R'
	Type Verb = 'T'
	Unmerged Verb = 'U'
	Unknown Verb = 'X'
)
func (v Verb) String() string {
	return string(v)
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
			out += color.String("\nRemoved "+ f.Src, color.Meta)
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
				older := strings.SplitAfter(string(chunk.Old), "\n",-1)
				newer := strings.SplitAfter(string(chunk.New), "\n",-1)
				lastline:=chunk.Line
				mychunks := patience.DiffFromLine(chunk.Line, older, newer)
				fmt.Println(color.String("¤¤¤¤¤¤ "+f.Src+string(chunk.Line)+" ¤¤¤¤¤¤", color.Meta))
				for _,ch := range mychunks {
					if ch.Line > lastline + 6 {
						for i:=lastline+1; i<lastline+4; i++ {
							fmt.Print(" ",newer[i-chunk.Line])
						}
						fmt.Println(color.String("¤¤¤¤¤¤ "+f.Src+string(chunk.Line-3)+" ¤¤¤¤¤¤", color.Meta))
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
		debug.Println("yeckgys")
	}
	return o
}

func genLsFilesE(args ...string) ([]string, os.Error) {
	o, e := git.Read("ls-files", args...)
	return splitOnNulls(o), e
}

func splitOnNulls(s string) []string {
	xs := strings.Split(s, "\000", -1)
	if len(xs[len(xs)-1]) == 0 {
		return xs[0:len(xs)-1]
	}
	return xs
}

var configList map[string]string

var once sync.Once

func ListConfig() map[string]string {
	once.Do(func() {
		o,_ := git.Read("config", "--null", "--list")
		configList = splitOnNullsAndLines(o)
	})
	return configList
}

func splitOnNullsAndLines(s string) (out map[string]string) {
	xs := strings.Split(s, "\000", -1)
	out = make(map[string]string)
	for _,x := range xs {
		kv := strings.Split(x, "\n", -1)
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

func Blob(b git.Hash) (o string, e os.Error) {
	return git.Read("cat-file", "blob", b.String())
}

func GetCommitHash(c git.Commitish) (git.CommitHash, os.Error) {
	if h,ok := c.(git.CommitHash); ok {
		return h,nil
	}
	return RevParse(c.String())
}

func rawCommit(c git.Commitish) (o string, e os.Error) {
	return git.Read("cat-file", "commit", c.String())
}

func Commit(c git.Commitish) (ce CommitEntry, e os.Error) {
	o,e := rawCommit(c)
	if e != nil { return }
	ls := strings.Split(o, "\n", -1)
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
		if n.String() != "0000000000000000000000000000000000000000" {
			args = stringslice.Append(args, n.String())
		}
	}
	for _,o := range older {
		if o.String() != "0000000000000000000000000000000000000000" {
			args = stringslice.Append(args, "^"+o.String())
		}
	}
	return RevListS(args)
}

func RevListS(args []string) (commits []git.CommitHash, e os.Error) {
	o,e := git.ReadS("rev-list", args)
	if e != nil { return }
	ls := strings.Split(o,"\n",-1)
	for _,l := range ls {
		if len(l) == 40 {
			commits = pars.Append(commits, git.CommitHash(mkHash(l)))
		}
	}
	return
}

func RevParse(rev string) (h git.CommitHash, e os.Error) {
	o,e := git.Read("rev-parse", "--verify", rev)
	if e != nil { return }
	return git.CommitHash(mkHash(o)), e
}

func MergeBase(us, them git.Commitish) (h git.CommitHash, e os.Error) {
	o, e := git.Read("merge-base", us.String(), them.String())
	if e != nil { return }
	return git.CommitHash(mkHash(o)), e
}

func MergeIndexAll() os.Error {
	return git.Run("merge-index", "git-merge-one-file", "-a")
}

func SendPack(repo0 string, updates map[git.Ref]git.CommitHash) os.Error {
	repo := RemoteUrl(repo0)
	args := []string{repo}
	for name,hash := range updates {
		args = stringslice.Append(args, hash.String()+":"+name.String())
	}
	return git.Run("send-pack", args...)
}
