package debug

import "fmt"
import "os"
import "goopt"

var amdebug = goopt.Bool("--debug", false, "enable debugging")

func Print(v ...interface{}) os.Error {
	if *amdebug {
		_,e := fmt.Fprintln(os.Stderr, v)
		return e
	}
	return nil
}
