package core

import (
	"os"
	"fmt"
	"strings"
	"bytes"
	"io"
	"io/ioutil"
	"../git/plumbing"
	"../git/color"
	"../util/out"
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

func (d *FileDiff) Fprint(f io.Writer) (e os.Error) {
	switch d.Change {
	case plumbing.Added:
		fmt.Fprint(f, color.String("Added "+d.Name, color.Meta))
	case plumbing.Deleted:
		fmt.Fprint(f,color.String("Deleted "+d.Name, color.Meta))
	case plumbing.Modified:
		if d.OldMode != d.NewMode { fmt.Fprintln(f,d) }
		oldf, e := plumbing.Blob(d.OldHash)
		if e != nil { return }
		newf, e := ioutil.ReadFile(d.Name)
		older := strings.SplitAfter(string(oldf),"\n",0)
		newer := strings.SplitAfter(string(newf),"\n",0)
		if e != nil { return }
		mychunks := patience.Diff(older, newer)
		
		chunkLine := 1
		lastline := chunkLine
		if len(mychunks) == 0 { return }
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

func (d *FileDiff) Show() string {
	f := bytes.NewBufferString("")
	d.Fprint(f)
	return f.String()
}

func (d *FileDiff) Print() os.Error {
	return d.Fprint(out.Writer)
}
