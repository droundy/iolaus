package main

import (
	"os"
	"fmt"
	"exec"
	"goopt"
	"path"
	"../src/util/debug"
	"../src/util/error"
	stringslice "./gotgo/slice(string)"
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
	pandoc("-o", *outname, mdf)
}

func announce(err os.Error) os.Error {
	debug.Print(err)
	return err
}

func pandoc(args ...string) (e os.Error) {
	debug.Print("calling pandoc",args)
	args = stringslice.Cat([]string{"pandoc"}, args)
	pandoc, e := exec.LookPath("pandoc")
	if e != nil { announce(e); return }
	pid,e := exec.Run(pandoc, args, os.Environ(), ".",
		exec.PassThrough, exec.PassThrough, exec.PassThrough)
	if e != nil { announce(e); return }
	ws,e := pid.Wait(0) // could have been os.WRUSAGE
	if e != nil { announce(e); return }
	if ws.ExitStatus() != 0 {
		e = os.NewError(fmt.Sprintf("pandoc exited with '%v'",ws.ExitStatus()))
		announce(e)
		return
	}
	return nil
}
