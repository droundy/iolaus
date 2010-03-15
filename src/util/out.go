package out

import (
	"fmt"
	"os"
	"io"
	"./cook"
)

var Writer io.Writer = os.Stdout

// Print to standard output.  Eventually I may add optional piping to
// a pager to this function, so it should be used in preference to
// fmt.Print itself.
func Print(v ...interface{}) os.Error {
	_,e := fmt.Fprint(Writer, v)
	return e
}

// Print to standard output.  Eventually I may add optional piping to
// a pager to this function, so it should be used in preference to
// fmt.Print itself.
func Println(v ...interface{}) os.Error {
	_,e := fmt.Fprintln(Writer, v)
	return e
}

func Printf(f string, v ...interface{}) os.Error {
	_,e := fmt.Fprintf(Writer, f, v)
	return e
}

func PromptForChar(f string, v ...interface{}) (byte, os.Error) {
	defer cook.Undo(cook.SetRaw())
	_,e := fmt.Printf(f, v)
	if e != nil { return 0, e }
	x := make([]byte,1)
	_,e = os.Stdin.Read(x)
	fmt.Print("\r\n")
	return x[0], e	
}
