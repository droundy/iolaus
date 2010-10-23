package porcelain

import (
	"./git"
	"os"
)

func Stash(args ...string) os.Error {
	return git.Run("stash", args...)
}

func Init() os.Error {
	return git.Run("init")
}
