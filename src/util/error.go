package error

import "fmt"
import "os"

func Print(v ...interface{}) os.Error {
	_,e := fmt.Fprintln(os.Stderr, v)
	return e
}

func FailOn(e os.Error) {
	if e != nil {
		Print(e)
		os.Exit(1)
	}
}

func Exit(e os.Error) {
	if e != nil {
		Print(e)
		os.Exit(1)
	}
	os.Exit(0)
}
