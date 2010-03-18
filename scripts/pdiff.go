package main;

import "fmt"
import "os"
import "strings"
import "io/ioutil"
import "../src/util/patience"

func main() {
	obytes, _ := ioutil.ReadFile(os.Args[1])
	nbytes, _ := ioutil.ReadFile(os.Args[2])
	o := strings.SplitAfter(string(obytes), "\n",0)
	n := strings.SplitAfter(string(nbytes), "\n",0)
	fmt.Println(patience.Diff(o,n))
}
