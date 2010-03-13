package main;

import (
	"./git/git"
	"./git/porcelain"
	"./util/error"
	"./util/help"
)

var description = func() string {
	return `
Call initialize once for each project you work on. Run it from the top
level directory of the project, with the project files already there.
Initialize will set up all the directories and files iolaus needs in order to
start keeping track of revisions for your project.
`}

func main() {
	help.Init("initialize a new repository.", description, nil)
	git.AmNotDirectlyInRepo("Cannot call iolaus-initialize in a repository!")
	error.Exit(porcelain.Init())
}
