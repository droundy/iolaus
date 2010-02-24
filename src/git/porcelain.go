package porcelain

import (
	"./git"
	"os"
)

func Init() os.Error {
	return git.Run("init")
}
