package help

import "fmt"
import "os"
import "flag"

func Init(summary string, getopts func() []string ) {
	h := flag.Bool("help",false,"show help message")
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of %s:\n\t", os.Args[0])
		fmt.Fprintln(os.Stderr, summary)
		flag.PrintDefaults()
	}
	listopts := flag.Bool("list-options",false,
		                    "list options in machine-readable format")
	flag.Parse()
	if *listopts {
		if getopts != nil {
			for _, o := range getopts() {
				fmt.Println(o)
			}
		}
		flag.VisitAll(func (f *flag.Flag) { fmt.Println("--"+f.Name) })
		os.Exit(0)
	}
	if *h {
		flag.Usage();
		os.Exit(0);
	}
}
