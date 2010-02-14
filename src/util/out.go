package out

import "fmt"
import "os"

// Print to standard output.  Eventually I may add optional piping to
// a pager to this function, so it should be used in preference to
// fmt.Println itself.
func Print(v ...interface{}) os.Error {
	_,e := fmt.Println(v)
	return e
}
