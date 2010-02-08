package plumbing

import (
	"./git"
	"os"
	"strings"
)

func Init() {
	git.Run("init",nil)
}

func LsFiles() (fs []string, e os.Error) {
	var o string
	o, e = git.Read("ls-files",
		              []string{"--exclude-standard", "-z", "--others", "--cached"})
	fs = strings.Split(o, "\000", 0)
	return
}
