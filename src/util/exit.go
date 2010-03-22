package exit

import (
	"os"
	hooks "./gotgo/slice(func())"
)

func Exit(ecode int) {
	pleaseExit <- ecode;
	// The following is to guarantee that we don't exit this function...
	AtExit(func() {})
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
var pleaseExit = make(chan int)
var aeRequests = make(chan aeReq)
func init() {
    go handleExit()
}

func handleExit() {
  atExit := []func(){}
	eraseOne := make(chan int)
	for {
		select {
		case x := <- aeRequests:
			i := len(atExit)
			atExit = hooks.Append(atExit, x.atexit)
			x.cancel <- func() { eraseOne <- i }
		case i := <- eraseOne:
			atExit[i] = nil
			for j := len(atExit)-1; j >= 0 && atExit[j] == nil; j-- {
				// clear buffer... this assumes that the cancel function is only
				// run once...
				atExit = atExit[0:j]
			}
		case ecode := <- pleaseExit:
			// no more exiting allowed! This should (hopefully) cause any
			// threads that try to AtExit to hang in a recognizable way (so
			// they may be GCed), and similarly any that try to Exit.
			for i:=len(atExit)-1; i>=0; i-- {
				if atExit[i] != nil {
					atExit[i]()
				}
			}
			os.Exit(ecode)
			return
		}
	}
}
