package git

import (
	"../util/debug"
	"os"
	"exec"
	"path"
	"fmt"
	"io/ioutil"
)

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
		wd = path.Clean(wd[0:len(wd)-1])
	}
	fmt.Println(mess)
	os.Exit(1)
}

func AmNotInRepo(mess string) {
	oldwd, _ := os.Getwd()
	for oldwd != "" {
		s, e := os.Stat(oldwd+"/.git")
		if e == nil && s != nil && s.IsDirectory() {
			fmt.Println(mess)
			os.Exit(1)
		}
		oldwd, _ = path.Split(oldwd)
		oldwd = path.Clean(oldwd[0:len(oldwd)-1])
	}
}

func AmNotDirectlyInRepo(mess string) {
	s, e := os.Stat(".git")
	if e == nil && s != nil && s.IsDirectory() {
		fmt.Println(mess)
		os.Exit(1)
	}
}

func announce(err os.Error) {
	debug.Print(err)
}

func Read(arg1 string, args []string) (output string, err os.Error) {
	debug.Print("calling git",arg1,args)
	args2 := make([]string, 2+len(args))
	args2[0] = "git" // zeroth argument is the program name...
	args2[1] = arg1 // first argument is the command...
	for i, c := range args {
    args2[i+2] = c
  }
	output = "" // empty output if we have an error...
	git, err := exec.LookPath("git")
	if err != nil { announce(err); return }
	pid,err := exec.Run(git, args2, os.Environ(),
		exec.PassThrough, exec.Pipe, exec.PassThrough)
	if err != nil { announce(err); return }
	o,err := ioutil.ReadAll(pid.Stdout)
	if err != nil { announce(err); return }
	ws,err := pid.Wait(0) // could have been os.WRUSAGE
	if err != nil { announce(err); return }
	if ws.ExitStatus() != 0 {
		err = os.NewError("git exited with "+string(ws.ExitStatus()))
		announce(err)
		return
	}
	return string(o), nil
}

func WriteRead(arg1 string, args []string, inp string) (output string, e os.Error) {
	debug.Print("calling git ",args)
	args2 := make([]string, 2+len(args))
	args2[0] = "git" // zeroth argument is the program name...
	args2[1] = arg1 // first argument is the command...
	for i, c := range args {
    args2[i+2] = c
  }
	output = "" // empty output if we have an error...
	git, e := exec.LookPath("git")
	if e != nil { announce(e); return }
	pid,e := exec.Run(git, args2, os.Environ(),
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
		e = os.NewError("git exited with "+string(ws.ExitStatus()))
		announce(e)
		return
	}
	return
}

func Write(arg1 string, args []string, inp string) (e os.Error) {
	debug.Print("calling git ",args)
	args2 := make([]string, 2+len(args))
	args2[0] = "git" // zeroth argument is the program name...
	args2[1] = arg1 // first argument is the command...
	for i, c := range args {
    args2[i+2] = c
  }
	git, e := exec.LookPath("git")
	if e != nil { announce(e); return }
	pid,e := exec.Run(git, args2, os.Environ(),
		exec.Pipe, exec.PassThrough, exec.PassThrough)
	if e != nil { announce(e); return }
	_,e = fmt.Fprint(pid.Stdin, inp)
	if e != nil { announce(e); return }
	e = pid.Stdin.Close()
	ws,e := pid.Wait(0) // could have been os.WRUSAGE
	if e != nil { announce(e); return }
	if ws.ExitStatus() != 0 {
		e = os.NewError("git exited with "+string(ws.ExitStatus()))
		announce(e)
		return
	}
	return nil
}

func Run(arg1 string, args []string) (e os.Error) {
	debug.Print("calling git ",args)
	args2 := make([]string, 2+len(args))
	args2[0] = "git" // zeroth argument is the program name...
	args2[1] = arg1 // first argument is the command...
	for i, c := range args {
    args2[i+2] = c
  }
	git, e := exec.LookPath("git")
	if e != nil { announce(e); return }
	pid,e := exec.Run(git, args2, os.Environ(),
		exec.PassThrough, exec.PassThrough, exec.PassThrough)
	if e != nil { announce(e); return }
	ws,e := pid.Wait(0) // could have been os.WRUSAGE
	if e != nil { announce(e); return }
	if ws.ExitStatus() != 0 {
		e = os.NewError("git exited with "+string(ws.ExitStatus()))
		announce(e)
		return
	}
	return nil
}
