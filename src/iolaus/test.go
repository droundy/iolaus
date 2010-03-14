package test

import (
	"os"
	"fmt"
	"exec"
	"strings"
	"syscall"
	"goopt"
	git "../git/git"
	"../git/plumbing"
	"../util/out"
	box "./gotgo/box(git.CommitHash,git.Commitish)"
)

var notest = goopt.Flag([]string{"--no-test"}, []string{"--test"},
	"do not run the test suite", "run the test suite [default]")

var machineName = func() string {
	var un syscall.Utsname
	syscall.Uname(&un)
	return charsToString(un.Nodename) + " (" + charsToString(un.Sysname) +
		" " + charsToString(un.Release) + ")"
}()

func charsToString(x [65]int8) string {
	n := ""
	for _,v := range x {
		if v == 0 { return n }
		n += string(byte(v))
	}
	return n
}

// test.Commit tests a given (presumably new) commit, and returns a
// modified commit that may contain a note indicating that it's been
// tested.
func Commit(h git.CommitHash) (outh git.CommitHash, e os.Error) {
	outh = h
	if *notest { return }
	t, e := plumbing.Commit(h)
	if e != nil { return }
	msg,e := Tree(t.Tree)
	if e != nil { return outh, e }
	newlog := t.Message
	if msg != "" {
		if strings.LastIndex(newlog, ":") <= strings.LastIndex(newlog, "\n") {
			newlog += "\n"
		}
		newlog += "\n" + msg
	}
	outh = plumbing.CommitTree(t.Tree, box.Box(t.Parents), newlog)
	return outh, e
}

// test.Tree tests a given tree, and returns an os.Error indicating
// whether the tests failed.
func Tree(h git.TreeHash) (msg string, e os.Error) {
	if *notest { return "",nil }
	_,e = plumbing.RevParse("refs/tested/"+h.String())
	// if there is no error, then we've already tested this tree, and
	// found it passed!
	if e == nil { return "",e }
	e = os.RemoveAll("/tmp/silly-testing")
	if e != nil { return "",e }
	e = os.MkdirAll("/tmp/silly-testing", 0777)
	if e != nil { return "",e }
	//here,e := os.Getwd()
	//if e != nil { return "",e }
	//e = os.Chdir("/tmp/silly-testing")
	//if e != nil { return "",e }
	//defer os.Chdir(here)
	plumbing.ReadTree(h, "--index-output=.git/index.tmp")
	// It's pretty hokey to use os.Setenv here rather than using exec to
	// set it directly, but it shouldn't be a problem as long as we
	// aren't calling git from multiple goroutines.
	e = os.Setenv("GIT_INDEX_FILE", ".git/index.tmp")
	if e != nil { return "",e }
	defer os.Setenv("GIT_INDEX_FILE", "")
	// don't forget the trailing slash in the prefix!
	e = plumbing.CheckoutIndex("-a", "--prefix=/tmp/silly-testing/")
	if e != nil { return "",e }
	// remove the temporary index, but don't worry if this fails...
	os.Remove(".git/index.tmp")
	bstat,e := os.Stat("/tmp/silly-testing/.build")
	if e == nil && (bstat.Permission() & 1) == 1 {
		out.Print("Running build!")
		// There is an executable .build, so run it!
		pid,e := exec.Run("/tmp/silly-testing/.build", []string{}, os.Environ(),
			"/tmp/silly-testing", exec.DevNull, exec.PassThrough, exec.MergeWithStdout)
		if e != nil { return "",e }
		ws,e := pid.Wait(0)
		if e != nil { return "",e }
		if ws.ExitStatus() != 0 {
			return "",os.NewError(fmt.Sprintf(".build exited with '%v'",ws.ExitStatus()))
		}
		msg = "Built-on: " + machineName
	}
	tstat,e := os.Stat("/tmp/silly-testing/.test")
	if e == nil && (tstat.Permission() & 1) == 1 {
		out.Print("Running test!")
		// There is an executable .test, so run it!
		pid,e := exec.Run("/tmp/silly-testing/.test", []string{}, os.Environ(),
			"/tmp/silly-testing", exec.DevNull, exec.PassThrough, exec.MergeWithStdout)
		if e != nil { return "", e }
		ws,e := pid.Wait(0)
		if e != nil { return "", e }
		if ws.ExitStatus() != 0 {
			return "",os.NewError(fmt.Sprintf(".test exited with '%v'",ws.ExitStatus()))
		}
		msg = "Tested-on: " + machineName
	}
	e = os.RemoveAll("/tmp/silly-testing")
	return
}

func isExecutable(f string) bool {
	stat,e := os.Stat(f)
	return e == nil && (stat.Permission() & 1) == 1
}
