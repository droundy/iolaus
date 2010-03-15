package core

import (
	"os"
	"../git/plumbing"
	stringslice "./gotgo/slice(string)"
)

func ModifiedFiles() []string {
	return stringslice.Cat(plumbing.DiffFilesModified([]string{}), plumbing.LsOthers())
}

func DiffFiles(paths []string) (ds []plumbing.FileDiff, e os.Error) {
	ds, e = plumbing.DiffFiles(paths)
	lends := len(ds)
	fs := plumbing.LsOthers()
	newds := make([]plumbing.FileDiff, lends+len(fs))
	for i,d := range ds {
		newds[i] = d
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
