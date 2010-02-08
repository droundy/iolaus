package main;

import "./git/git"
import "./git/plumbing"
import "./util/error"
import "./util/help"

func main() {
	git.AmInRepo("Must be in a repository to call whatsnew!")
	help.Init("see unrecorded changes.")
	fs,_ := plumbing.LsFiles()
	error.Print(fs)
	error.Print("This isn't working yet, obviously...")
}
