package plumbing

import (
	"./git"
	"os"
	"strings"
)

func Init() {
	git.Run("init",nil)
}

func LsFiles() []string {
	o, _ := LsFilesE()
	return o
}

func LsFilesE() (fs []string, e os.Error) {
	var o string
	o, e = git.Read("ls-files",
		              []string{"--exclude-standard", "-z", "--others", "--cached"})
	fs = strings.Split(o, "\000", 0)
	return fs[0:len(fs)-1], e
}
