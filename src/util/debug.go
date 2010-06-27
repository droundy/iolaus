package debug

import "fmt"
import "os"
import "github.com/droundy/goopt"

var amdebug = goopt.Bool("--debug", false, "enable debugging")

func Print(v ...interface{}) os.Error {
	if *amdebug {
		_,e := fmt.Fprint(os.Stderr, v)
		return e
	}
	return nil
}

func Printf(f string, v ...interface{}) os.Error {
	if *amdebug {
		_,e := fmt.Fprintf(os.Stderr, f, v)
		return e
	}
	return nil
}

func Println(v ...interface{}) os.Error {
	if *amdebug {
		_,e := fmt.Fprintln(os.Stderr, v)
		return e
	}
	return nil
}
