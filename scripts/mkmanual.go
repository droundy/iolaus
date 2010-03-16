package main

import (
	"os"
	"fmt"
	"exec"
	"io/ioutil"
	"../src/util/debug"
	"../src/util/error"
	stringslice "./gotgo/slice(string)"
)

func main() {
	os.MkdirAll("doc", 0777)
	header,err := ioutil.ReadFile("scripts/header.html")
	error.FailOn(err)
	html := string(header)
	html += "<ul>\n"
	for _,cmd := range os.Args[1:] {
		cmd = cmd[4:len(cmd)-3]
		html += "<li><a href=\""+cmd+".html\">"+cmd+"</a></li>\n"
	}
	html += "</ul>\n"
	footer,err := ioutil.ReadFile("scripts/footer.html")
	error.FailOn(err)
	ioutil.WriteFile("doc/manual.html", []byte(html+string(footer)), 0666)
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
