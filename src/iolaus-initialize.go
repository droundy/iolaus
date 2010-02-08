package main;

import "./git/git"
import "./git/porcelain"
import "./util/error"
import "./util/help"

func main() {
	git.AmNotDirectlyInRepo("Cannot call iolaus-initialize in a repository!")
	help.Init("initialize a new repository.", nil)
	error.Exit(porcelain.Init())
}
