package main;

import "./git/git"
import "./git/plumbing"
import "./util/error"
import "./util/help"

func main() {
	git.AmInRepo("Must be in a repository to call whatsnew!")
	help.Init("see unrecorded changes.", plumbing.LsFiles)
	error.Print("This isn't working yet, obviously...")
}
