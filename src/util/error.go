package error

import (
	"fmt"
	"os"
	"./exit"
	"./cook"
)

func Print(v ...interface{}) os.Error {
	_,e := fmt.Fprint(os.Stderr, v)
	fmt.Fprint(os.Stderr, "\r\n")
	return e
}

func FailOn(e os.Error) {
	if e != nil {
		cook.SetCooked()
		Print(e,"\n")
		exit.Exit(1)
	}
}

func Exit(e os.Error) {
	cook.SetCooked()
	if e != nil {
		Print(e,"\n")
		exit.Exit(1)
	}
	exit.Exit(0)
}
