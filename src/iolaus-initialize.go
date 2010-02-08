package main;

import "./git/git"
import "./git/porcelain"
import "./util/error"
import "./util/help"

func main() {
	help.Init("initialize a new repository.")
	git.AmNotDirectlyInRepo("Cannot call iolaus-initialize in a repository!")
	error.Exit(porcelain.Init())
}
