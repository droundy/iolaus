package help

import (
	"fmt"
	"os"
	"goopt"
)

func Init(summary string, getopts func() []string ) {
	goopt.Usage = func() string {
		return fmt.Sprintf("Usage of %s:\n\t",os.Args[0]) +
			summary + "\n" + goopt.Help()
	}
	listopts := func () os.Error {
		if getopts != nil {
			for _, o := range getopts() {
				fmt.Println(o)
			}
		}
		goopt.VisitAllNames(func (n string) { fmt.Println(n) })
		os.Exit(0)
		return nil
	}
	goopt.NoArg([]string{"--list-options"},
		"list options in machine-readable format", listopts)
	goopt.Parse()
}
