package core

import (
	"os"
	"fmt"
	"strings"
	"bytes"
	"io"
	"io/ioutil"
	"../git/git"
	"../git/plumbing"
	"../git/color"
	"../util/out"
	"../util/debug"
	"../util/patience"
	stringslice "./gotgo/slice(string)"
)

func ModifiedFiles() []string {
	return stringslice.Cat(plumbing.DiffFilesModified([]string{}), plumbing.LsOthers())
}

func DiffFiles(paths []string) (ds []FileDiff, e os.Error) {
	oldds, e := plumbing.DiffFiles(paths)
	lends := len(oldds)
	fs := plumbing.LsOthers()
	newds := make([]FileDiff, lends+len(fs))
	for i,d := range oldds {
		newds[i] = FileDiff(d)
	}
	ds = newds
	for i,f := range fs {
		ds[lends+i].Change = plumbing.Added
		ds[lends+i].OldMode = 0
		ds[lends+i].NewMode = 0
		ds[lends+i].Name = f
		for j := range ds[lends+i].OldHash {
			ds[lends+i].OldHash[j] = '0'
			ds[lends+i].NewHash[j] = '0'
		}
	}
	return
}

type FileDiff plumbing.FileDiff

func (d *FileDiff) UpdateNew(contents string) (e os.Error) {
	switch d.Change {
	case plumbing.Modified:
		h,e := plumbing.HashObject(contents)
		if e != nil { return }
		d.NewHash = h
	}
	return
}

func (d *FileDiff) OldNewStrings() (o, n string, e os.Error) {
	switch d.Change {
	default:
		oldf, e := plumbing.Blob(d.OldHash)
		if e != nil { return }
		o = string(oldf)
		if d.NewHash.IsEmpty() {
			newf, err := ioutil.ReadFile(d.Name)
			if err != nil {
				e = err
				return
			}
			n = string(newf)
		} else {
			n, e = plumbing.Blob(d.NewHash)
		}
	}
	return
}

func (d *FileDiff) Fprint(f io.Writer) (e os.Error) {
	switch d.Change {
	case plumbing.Added:
		fmt.Fprintln(f, color.String("Added "+d.Name, color.Meta))
	case plumbing.Deleted:
		fmt.Fprintln(f,color.String("Deleted "+d.Name, color.Meta))
	case plumbing.Modified:
		if d.OldMode != d.NewMode { fmt.Fprintln(f,d) }
		oldf, newf, e := d.OldNewStrings()
		if e != nil {
			debug.Printf("Had trouble reading file and old %s: %s\n", d.Name, e.String())
			return
		}
		if newf == oldf {
			debug.Println("File "+d.Name+" is unchanged.")
			return
		}
		newer := strings.SplitAfter(newf,"\n",0)
		older := strings.SplitAfter(oldf,"\n",0)
		mychunks := patience.Diff(older, newer)
		
		chunkLine := 1
		lastline := chunkLine
		debug.Printf("File %s has %d chunks changed\n",d.Name,len(mychunks))
		if len(mychunks) == 0 {
			debug.Println("File "+d.Name+" mysteriously looks unchanged.")
			return
		}
		if mychunks[0].Line-4 > lastline {
			lastline = mychunks[0].Line - 4
		}
		fmt.Fprintf(f,color.String("¤¤¤ %s %d ¤¤¤\n", color.Meta),d.Name,lastline)
		for _,ch := range mychunks {
			if ch.Line > lastline + 6 {
				for i:=lastline+1; i<lastline+4; i++ {
					fmt.Fprint(f," ",newer[i-chunkLine])
				}
				fmt.Fprintf(f,color.String("¤¤¤ %s %d ¤¤¤\n", color.Meta),d.Name,ch.Line-3)
				for i:=ch.Line-3; i<ch.Line; i++ {
					fmt.Fprint(f," ",newer[i-chunkLine])
				}
			} else {
				for i:=lastline; i<ch.Line; i++ {
					fmt.Fprint(f," ",newer[i-chunkLine])
				}
			}
			fmt.Fprint(f,ch)
			lastline = ch.Line + len(ch.New)
		}
		for i:=lastline-chunkLine; i<len(newer)-1 && i < lastline-chunkLine+3;i++ {
			fmt.Fprint(f," ",newer[i])
		}
	default:
		fmt.Fprintln(f,d)
	}
	return
}

// Check if the file has something changed!
func (d *FileDiff) HasChange() bool {
	switch d.Change {
	case plumbing.Added, plumbing.Deleted:
		return true
	case plumbing.Modified:
		if d.OldMode != d.NewMode { return true }
		oldf, newf, e := d.OldNewStrings()
		if e != nil { return true }
		debug.Printf("File %s has changes? %v\n", d.Name, oldf != newf)
		return oldf != newf
	}
	return true
}

func (d *FileDiff) Show() string {
	f := bytes.NewBufferString("")
	d.Fprint(f)
	return f.String()
}

func (d *FileDiff) Print() os.Error {
	return d.Fprint(out.Writer)
}

func isEmpty(h git.Hash) bool {
	for _,v := range h {
		if v != '0' { return false }
	}
	return true
}

func (d *FileDiff) Info() (mode int, contents git.Hash, f string) {
	mode = d.NewMode
	contents = d.NewHash
	if isEmpty(contents) && d.Change != plumbing.Deleted {
		contents,_ = plumbing.HashFile(d.Name)
	}
	return mode, contents, d.Name
}

// Merge two commits, and return the resulting tree.
func Merge(us, them git.Commitish) (t git.TreeHash, e os.Error) {
	var basetree git.TreeHash
	base, e := plumbing.MergeBase(us, them)
	if e == nil {
		baseX,e := plumbing.Commit(base)
		if e != nil { return }
		basetree = baseX.Tree
	} else {
		debug.Println("Looks like there are no common ancestors.")
		// In this case, we want to use the empty tree for the merge.
	}
	e = plumbing.ReadTree3(basetree, us, them) 
	if e != nil { return }
	e = plumbing.MergeIndexAll()
	if e != nil { return }
	return plumbing.WriteTree(), e
}
