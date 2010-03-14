package core

import (
	"../git/plumbing"
	stringslice "./gotgo/slice(string)"
)

func ModifiedFiles() []string {
	return stringslice.Cat(plumbing.DiffFilesModified([]string{}), plumbing.LsOthers())
}
