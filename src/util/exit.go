package exit

import (
	"os"
)

func Exit(ecode int) {
	ae := make(chan func())
	pleaseExit <- ae;
	atexit := <- ae
	atexit()
	os.Exit(ecode)
}

// returns a "cancel this AtExit" function
func AtExit(f func()) func() {
	canc := make(chan func())
	aeRequests <- aeReq{ canc, f }
	return <- canc
}

type aeReq struct {
	cancel chan <- func()
	atexit func()
}
var pleaseExit = make(chan (chan <- func()))
var aeRequests = make(chan aeReq)
func init() {
    go handleExit()
}

func append(fs []func(), f func()) []func() {
	if len(fs) == cap(fs) {
		newfs := make([]func(), len(fs), 2*(len(fs)+1))
		for i,v := range(fs) { newfs[i] = v }
		fs = newfs
	}
	fs = fs[0:len(fs)+1]
	fs[len(fs)-1] = f
	return fs
}

func handleExit() {
  atExit := make([]func(), 0, 100)
	eraseOne := make(chan int)
	for {
		select {
		case x := <- aeRequests:
			i := len(atExit)
			atExit = append(atExit, x.atexit)
			x.cancel <- func() { eraseOne <- i }
		case i := <- eraseOne:
			atExit[i] = nil
			for j := len(atExit)-1; j >= 0 && atExit[j] == nil; j-- {
				// clear buffer... this assumes that the cancel function is only
				// run once...
				atExit = atExit[0:j]
			}
		case ch := <- pleaseExit:
			f := func() {
				for i:=len(atExit)-1; i>=0; i-- { if atExit[i] != nil { atExit[i]() } }
			}
			ch <- f
			// no more exiting allowed! This should (hopefully) cause any
			// threads that try to AtExit to hang in a recognizable way (so
			// they may be GCed), and similarly any that try to Exit.
			return
		}
	}
}
