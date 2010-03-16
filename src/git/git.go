package git

import (
	"../util/debug"
	"../util/exit"
	"os"
	"exec"
	"path"
	"fmt"
	"io/ioutil"
	stringslice "./gotgo/slice(string)"
)

// These data types are mostly just defined here so that I can use
// gotgo on them conveniently...
type Hash [40]byte
func (h Hash) String() string { return string(h[0:40]) }
func (h Hash) IsEmpty() bool {
	for _,v := range h {
		if v != '0' { return false }
	}
	return true
}
type TreeHash Hash
func (r TreeHash) String() string { return string(r[0:40]) }
func (r TreeHash) treeString() string { return string(r[0:40]) }
type Ref string
func (r Ref) String() string { return string(r) }
func (r Ref) commitString() string { return string(r) }
func (r Ref) treeString() string { return string(r) }
type CommitHash Hash
func (h CommitHash) String() string { return string(h[0:40]) }
func (h CommitHash) commitString() string { return string(h[0:40]) }
func (h CommitHash) treeString() string { return string(h[0:40]) }
type Treeish interface {
	String() string
	treeString() string
}
type Commitish interface {
	String() string
	commitString() string
}

func AmInRepo(mess string) {
	oldwd, _ := os.Getwd()
	wd := oldwd
	for wd != "" {
		s, e := os.Stat(wd+"/.git")
		if e == nil && s != nil && s.IsDirectory() {
			os.Chdir(wd)
			return
		}
		wd, _ = path.Split(wd)
		if wd != "" {
			wd = path.Clean(wd[0:len(wd)-1])
		}
	}
	fmt.Println(mess)
	exit.Exit(1)
}

func AmNotInRepo(mess string) {
	oldwd, _ := os.Getwd()
	for oldwd != "" {
		s, e := os.Stat(oldwd+"/.git")
		if e == nil && s != nil && s.IsDirectory() {
			fmt.Println(mess)
			exit.Exit(1)
		}
		oldwd, _ = path.Split(oldwd)
		oldwd = path.Clean(oldwd[0:len(oldwd)-1])
	}
}

func AmNotDirectlyInRepo(mess string) {
	s, e := os.Stat(".git")
	if e == nil && s != nil && s.IsDirectory() {
		fmt.Println(mess)
		exit.Exit(1)
	}
}

func announce(err os.Error) os.Error {
	debug.Println(err)
	return err
}

func explain(s string, err os.Error) os.Error {
	return announce(os.NewError(s+": "+err.String()))
}

func Read(arg1 string, args ...string) (output string, err os.Error) {
	return ReadS(arg1, args)
}

func WriteRead(arg1 string, inp string, args ...string) (output string, e os.Error) {
	return WriteReadS(arg1, inp, args)
}

func Write(arg1 string, inp string, args ...string) (e os.Error) {
	return WriteS(arg1, inp, args)
}

func Run(arg1 string, args ...string) (e os.Error) {
	return RunS(arg1, args)
}

func RunSilently(arg1 string, args ...string) (e os.Error) {
	return RunSilentlyS(arg1, args)
}

func ReadS(arg1 string, args []string) (output string, err os.Error) {
	debug.Println("calling git",arg1,args)
	args = stringslice.Cat([]string{"git", arg1}, args)
	output = "" // empty output if we have an error...
	git, err := exec.LookPath("git")
	if err != nil { err = explain("exec.LookPath",err); return }
	pid,err := exec.Run(git, args, os.Environ(), ".",
		exec.PassThrough, exec.Pipe, exec.PassThrough)
	if err != nil { announce(err); return }
	o,err := ioutil.ReadAll(pid.Stdout)
	if err != nil { announce(err); return }
	ws,err := pid.Wait(0) // could have been os.WRUSAGE
	if err != nil { announce(err); return }
	if ws.ExitStatus() != 0 {
		err = os.NewError(fmt.Sprintf("git %s exited with '%v'",arg1,ws.ExitStatus()))
		announce(err)
		return
	}
	return string(o), nil
}

func WriteReadS(arg1 string, inp string, args []string) (output string, e os.Error) {
	debug.Println("calling git",arg1,args)
	args = stringslice.Cat([]string{"git", arg1}, args)
	output = "" // empty output if we have an error...
	git, e := exec.LookPath("git")
	if e != nil { e = explain("exec.LookPath",e); return }
	pid,e := exec.Run(git, args, os.Environ(), ".",
		exec.Pipe, exec.Pipe, exec.PassThrough)
	if e != nil { announce(e); return }
	_,e = fmt.Fprint(pid.Stdin, inp)
	if e != nil { announce(e); return }
	e = pid.Stdin.Close()
	if e != nil { announce(e); return }
	o,e := ioutil.ReadAll(pid.Stdout)
	output = string(o)
	if e != nil { announce(e); return }
	ws,e := pid.Wait(0) // could have been os.WRUSAGE
	if e != nil { announce(e); return }
	if ws.ExitStatus() != 0 {
		e = os.NewError(fmt.Sprintf("git exited with '%v'",ws.ExitStatus()))
		announce(e)
		return
	}
	return
}

func WriteS(arg1 string, inp string, args []string) (e os.Error) {
	debug.Println("calling git",arg1,args)
	args = stringslice.Cat([]string{"git", arg1}, args)
	git, e := exec.LookPath("git")
	if e != nil { announce(e); return }
	pid,e := exec.Run(git, args, os.Environ(), ".",
		exec.Pipe, exec.PassThrough, exec.PassThrough)
	if e != nil { announce(e); return }
	_,e = fmt.Fprint(pid.Stdin, inp)
	if e != nil { announce(e); return }
	e = pid.Stdin.Close()
	ws,e := pid.Wait(0) // could have been os.WRUSAGE
	if e != nil { announce(e); return }
	if ws.ExitStatus() != 0 {
		e = os.NewError(fmt.Sprintf("git exited with '%v'",ws.ExitStatus()))
		announce(e)
		return
	}
	return nil
}

func RunS(arg1 string, args []string) (e os.Error) {
	debug.Println("calling git",arg1,args)
	args = stringslice.Cat([]string{"git", arg1}, args)
	git, e := exec.LookPath("git")
	if e != nil { announce(e); return }
	pid,e := exec.Run(git, args, os.Environ(), ".",
		exec.PassThrough, exec.PassThrough, exec.PassThrough)
	if e != nil { announce(e); return }
	ws,e := pid.Wait(0) // could have been os.WRUSAGE
	if e != nil { announce(e); return }
	if ws.ExitStatus() != 0 {
		e = os.NewError(fmt.Sprintf("git exited with '%v'",ws.ExitStatus()))
		announce(e)
		return
	}
	return nil
}

func RunSilentlyS(arg1 string, args []string) (e os.Error) {
	debug.Println("calling git",arg1,args)
	args = stringslice.Cat([]string{"git", arg1}, args)
	git, e := exec.LookPath("git")
	if e != nil { announce(e); return }
	pid,e := exec.Run(git, args, os.Environ(), ".",
		exec.PassThrough, exec.DevNull, exec.DevNull)
	if e != nil { announce(e); return }
	ws,e := pid.Wait(0) // could have been os.WRUSAGE
	if e != nil { announce(e); return }
	if ws.ExitStatus() != 0 {
		e = os.NewError(fmt.Sprintf("git exited with '%v'",ws.ExitStatus()))
		announce(e)
		return
	}
	return nil
}
