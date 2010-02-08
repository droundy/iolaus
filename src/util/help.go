package help

import "fmt"
import "os"
import "flag"

func Init(summary string) {
	h := flag.Bool("help",false,"show help message")
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of %s:\n\t", os.Args[0])
		fmt.Fprintln(os.Stderr, summary)
		flag.PrintDefaults()
	}
	flag.Parse()
	if *h {
		flag.Usage();
		os.Exit(0);
	}
}
