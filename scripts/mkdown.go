package main

import (
	"os"
	"fmt"
	"exec"
	"github.com/droundy/goopt.git"
	"path"
	"io/ioutil"
	"../src/util/debug"
	"../src/util/error"
	stringslice "../src/util/slice(string)"
)

var outname = goopt.String([]string{"-o","--output"}, "FILENAME",
	"name of output file")

func main() {
	goopt.Parse(nil)
	if len(goopt.Args) != 2 {
		error.Exit(os.NewError(os.Args[0]+" requires just one argument!"))
	}
	mdf := goopt.Args[1]
	if mdf[len(mdf)-3:] != ".md" {
		error.Exit(os.NewError(mdf+" doesn't end with .md"))
	}
	basename := mdf[0:len(mdf)-3]
	if *outname == "FILENAME" {
		*outname = basename+".html"
	}
	dir,_ := path.Split(*outname)
	os.MkdirAll(dir, 0777)
	body,err := pandoc(mdf)
	error.FailOn(err)
	header,err := ioutil.ReadFile("scripts/header.html")
	error.FailOn(err)
	footer,err := ioutil.ReadFile("scripts/footer.html")
	error.FailOn(err)
	ioutil.WriteFile(*outname, []byte(string(header)+body+string(footer)), 0666)
	error.Exit(nil)
}

func announce(err os.Error) os.Error {
	debug.Println(err)
	return err
}

func pandoc(args ...string) (out string, e os.Error) {
	debug.Println("calling pandoc",args)
	args = stringslice.Cat([]string{"pandoc"}, args)
	pandoc, e := exec.LookPath("pandoc")
	if e != nil { announce(e); return }
	pid,e := exec.Run(pandoc, args, os.Environ(), ".",
		exec.PassThrough, exec.Pipe, exec.PassThrough)
	if e != nil { announce(e); return }
	o,err := ioutil.ReadAll(pid.Stdout)
	if err != nil { announce(err); return }
	ws,e := pid.Wait(0) // could have been os.WRUSAGE
	if e != nil { announce(e); return }
	if ws.ExitStatus() != 0 {
		e = os.NewError(fmt.Sprintf("pandoc exited with '%v'",ws.ExitStatus()))
		announce(e)
		return
	}
	return string(o),nil
}
