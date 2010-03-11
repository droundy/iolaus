package main;

import (
	"./git/git"
	"./git/porcelain"
	"./util/error"
	"./util/help"
)

func main() {
	help.Init("initialize a new repository.", nil)
	git.AmNotDirectlyInRepo("Cannot call iolaus-initialize in a repository!")
	error.Exit(porcelain.Init())
}
