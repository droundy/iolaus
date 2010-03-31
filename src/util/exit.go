package exit

import (
	"os"
)

func Exit(ecode int) {
	pleaseExit <- struct{}{};
	_ = <- finishedAtExit
	os.Exit(ecode)
}

// returns a "cancel this AtExit" function
func AtExit(f func()) {
	aeRequests <- f
}

var finishedAtExit = make(chan struct{})
var pleaseExit = make(chan struct{})
var aeRequests = make(chan func ())
func init() {
    go handleExit()
}

func handleExit() {
	defer func() { finishedAtExit <- struct{}{} }()
	for {
		select {
		case x := <- aeRequests:
			defer x()
		case _ = <- pleaseExit:
			// no more exiting allowed! This should (hopefully) cause any
			// threads that try to AtExit to hang in a recognizable way (so
			// they may be GCed), and similarly any that try to Exit.
			return
		}
	}
}
