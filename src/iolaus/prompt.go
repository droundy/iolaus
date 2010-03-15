package prompt

import (
	"goopt"
	"./core"
	"../util/out"
	"../util/error"
	"../util/debug"
)

var all = goopt.Flag([]string{"-a","--all"}, []string{"--interactive"},
	"use all patches", "prompt for patches interactively")

func Run(ds []core.FileDiff, f func(core.FileDiff)) {
	if *all {
		for _,d := range ds {
			f(d)
		}
	} else {
	  files: for _,d := range ds {
			for {
				// Just keep asking until we get a reasonable answer...
				c,e := out.PromptForChar("Record changes to %s? ", d.Name)
				error.FailOn(e)
				switch c {
				case 'q','Q': error.Exit(e)
		    case 'v','V':
					d.Print()
				case 'y','Y':
					debug.Println("Dealing with file ",d.Name)
					f(d)
					continue files
				case 'n','N': debug.Println("Ignoring changes to file ",d.Name)
					continue files
				}
			}
		}
	}
}
