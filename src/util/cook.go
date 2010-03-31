package cook

import (
	"exec"
	"os"
	"once"
	"io/ioutil"
	"./exit"
)

// The secret type only exists so noone can use Undo to do anything
// but undo one of the below... which just goes to prove that I can't
// quit using type witnesses.
type secret string

func Undo(state secret) {
	if state == secret("") { return } // nothing to undo
	pid,e := exec.Run("/bin/stty", []string{"/bin/stty",string(state)},
		os.Environ(), ".",
		exec.PassThrough, exec.PassThrough, exec.DevNull)
	if e != nil { return }
	pid.Wait(0)
}

func readStty() secret {
	pid,e := exec.Run("/bin/stty", []string{"/bin/stty","-g"}, os.Environ(),
		".", exec.PassThrough, exec.Pipe, exec.DevNull)
	if e != nil { return secret("") }
	o,e := ioutil.ReadAll(pid.Stdout)
	if e != nil { return secret("") }
	ws,e := pid.Wait(0)
	if e != nil { return secret("") }
	if ws.ExitStatus() != 0 {
		return secret("")
		//panic(os.NewError("stty exited with "+string(ws.ExitStatus())))
	}
	return secret(o[0:len(o)-1])
}

func SetRaw() secret {
	once.Do(func() { exit.AtExit(func () { SetCooked() })})
	x := readStty() // could use "-echo" below...
	pid,e := exec.Run("/bin/stty", []string{"/bin/stty","raw"},
		os.Environ(), ".",
		exec.PassThrough, exec.PassThrough, exec.DevNull)
	if e != nil { return secret("") }
	pid.Wait(0)
	return x
}

func SetCooked() secret {
	x := readStty()
	pid,e := exec.Run("/bin/stty", []string{"/bin/stty","cooked","echo"},
		os.Environ(), ".",
		exec.PassThrough, exec.PassThrough, exec.DevNull)
	if e != nil { panic(e) }
	pid.Wait(0)
	return x
}
