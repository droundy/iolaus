package prompt

import (
	"goopt"
	"./core"
	"../util/out"
	"../util/error"
)

var All = goopt.Flag([]string{"-a","--all"}, []string{"--interactive"},
	"verb all patches", "prompt for patches interactively")

func Run(ds []core.FileDiff, f func(core.FileDiff)) {
	if *All {
		for _,d := range ds {
			f(d)
		}
	} else {
	  files: for _,d := range ds {
			for {
				// Just keep asking until we get a reasonable answer...
				c,e := out.PromptForChar(goopt.Expand("Verb changes to %s? "), d.Name)
				error.FailOn(e)
				switch c {
				case 'q','Q': error.Exit(e)
		    case 'v','V':
					d.Print()
				case 'y','Y':
					out.Println("Dealing with file ",d.Name)
					f(d)
					continue files
				case 'n','N': out.Println("Ignoring changes to file ",d.Name)
					continue files
				}
			}
		}
	}
}
