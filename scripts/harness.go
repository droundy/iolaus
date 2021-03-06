package main

import (
	"os"
	"fmt"
	"path"
	"exec"
	"io/ioutil"
	"../src/util/error"
	"../src/util/exit"
)

func main() {
	ts, e := os.Open("tests", os.O_RDONLY, 0)
	error.FailOn(e)
	error.FailOn(os.RemoveAll("tests/tmp"))
	ds, e := ts.Readdir(-1)
	error.FailOn(e)
	// Set up an environment with an appropriate path
	wd, e := os.Getwd()
	error.FailOn(e)
	os.Setenv("PATH", path.Join(wd,"bin")+":"+os.Getenv("PATH"))
	for _,d := range ds {
		if d.IsRegular() && (d.Permission() & 1) == 1 && endsWith(d.Name, ".sh") {
			basename := d.Name[0:len(d.Name)-3]
			fmt.Printf("Running %s... ", basename)
			dirname := path.Join("tests/tmp", basename)
			error.FailOn(os.MkdirAll(dirname, 0777))
			pid, e := exec.Run(path.Join(wd,"tests",d.Name), []string{}, os.Environ(), dirname,
				exec.DevNull, exec.Pipe, exec.MergeWithStdout)
			error.FailOn(e)
			o,e := ioutil.ReadAll(pid.Stdout)
			error.FailOn(e)
			ret, e := pid.Wait(0)
			error.FailOn(e)
			if ret.ExitStatus() != 0 {
				error.Print("FAILED!\n", string(o))
				error.Print("Test '", basename, "' failed!")
				exit.Exit(1)
			}
			fmt.Println("passed.")
		}
	}
}

func endsWith(s string, e string) bool {
	return len(s) >= len(e) && s[len(s)-len(e):] == e
}

func startsWith(s string, e string) bool {
	return len(s) >= len(e) && s[0:len(e)] == e
}
