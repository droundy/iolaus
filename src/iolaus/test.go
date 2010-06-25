package test

import (
	"os"
	"fmt"
	"exec"
	"path"
	"strings"
	"syscall"
	"github.com/droundy/goopt.git"
	git "../git/git"
	"../git/plumbing"
	"../util/out"
	"../util/debug"
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
		if strings.LastIndex(newlog, ":") <= strings.LastIndex(newlog, "\n") &&
			newlog[len(newlog)-1] != '\n' {
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
	testdir,e := tmpDir("/tmp/silly-testing")
	defer os.RemoveAll(testdir) // FIXME: this should be optional, eventually.
	if e != nil { return "",e }
	//here,e := os.Getwd()
	//if e != nil { return "",e }
	//e = os.Chdir(testdir)
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
	e = plumbing.CheckoutIndex("-a", "--prefix="+testdir+"/")
	if e != nil { return "",e }
	// remove the temporary index, but don't worry if this fails...
	os.Remove(".git/index.tmp")
	bstat,e := os.Stat(path.Join(testdir,".build"))
	if e == nil && (bstat.Permission() & 1) == 1 {
		out.Print("Running build!")
		// There is an executable .build, so run it!
		// clear out the GIT_INDEX_FILE environment variable, so our test
		// will be clean.
		pid,e := exec.Run(path.Join(testdir,".build"), []string{},
			removeEnv("GIT_INDEX_FILE", os.Environ()),
			testdir, exec.DevNull, exec.PassThrough, exec.MergeWithStdout)
		if e != nil { return "",e }
		ws,e := pid.Wait(0)
		if e != nil { return "",e }
		if ws.ExitStatus() != 0 {
			return "",os.NewError(fmt.Sprintf(".build exited with '%v'",ws.ExitStatus()))
		}
		msg = "Built-on: " + machineName
	}
	tstat,e := os.Stat(path.Join(testdir,".test"))
	if e == nil && (tstat.Permission() & 1) == 1 {
		out.Print("Running test!")
		// There is an executable .test, so run it!
		pid,e := exec.Run(path.Join(testdir,".test"), []string{},
			removeEnv("GIT_INDEX_FILE", os.Environ()),
			testdir, exec.DevNull, exec.PassThrough, exec.MergeWithStdout)
		if e != nil { return "", e }
		ws,e := pid.Wait(0)
		if e != nil { return "", e }
		if ws.ExitStatus() != 0 {
			return "",os.NewError(fmt.Sprintf(".test exited with '%v'",ws.ExitStatus()))
		}
		msg = "Tested-on: " + machineName
	} else if e == nil {
		debug.Printf("The test isn't executable... %o\n", tstat.Permission())
	} else {
		debug.Printf("There is no test... %v\n", e)
	}
	e = os.RemoveAll(testdir)
	return
}

func isExecutable(f string) bool {
	stat,e := os.Stat(f)
	return e == nil && (stat.Permission() & 1) == 1
}

func tmpDir(p string) (string, os.Error) {
	e := os.Mkdir(p,0777)
	if e == nil { return p,nil }
	for i:=0; i<30; i++ {
		pnew := p+"-"+fmt.Sprint(i)
		e = os.Mkdir(pnew,0777)
		if e == nil { return pnew,nil }
	}
	return "",e
}

func removeEnv(torem string, e []string) []string {
	out := make([]string,len(e))
	torem = torem + "="
	off := 0
	for i,v := range e {
		if len(v) > len(torem) && v[0:len(torem)] == torem {
			off -= 1
		} else {
			out[i+off] = v
		}
	}
	return out[0:len(out)+off]
}
