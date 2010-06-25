package prompt

import (
	"github.com/droundy/goopt.git"
	"strings"
	"./core"
	"../git/color"
	"../util/out"
	"../util/error"
	"../util/debug"
	"../util/patience"
	ss "../util/slice(string)"
)

// To make --all the default, set *prompt.All to true.
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
				if !d.HasChange() {
					continue files
				}
				// Just keep asking until we get a reasonable answer...
				c,e := out.PromptForChar(goopt.Expand("Verb changes to %s? "), d.Name)
				error.FailOn(e)
				switch c {
				case 'q','Q': error.Exit(e)
		    case 'v','V':
					d.Print()
		    case 'z','Z':
					o,n,e := d.OldNewStrings()
					error.FailOn(e)
					d.UpdateNew(pickChanges(d.Name, o, n))
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

func pickChanges(filename, o, n string) string {
	older := strings.SplitAfter(o,"\n",0)
	outer := strings.SplitAfter(o,"\n",0)
	newer := strings.SplitAfter(n,"\n",0)
	mychunks := patience.Diff(older, newer)
	offset := 0
	
chunks: for _,ch := range mychunks {
		// Here we deal with chunk "ch"
		for {
			// Just keep asking until we get a reasonable answer...
			lin := 0
			if ch.Line > 4 {
				lin = ch.Line - 4
			}
			out.Printf(color.String("¤¤¤ %s %d ¤¤¤\n", color.Meta),filename,offset+lin+1)
			for lin=lin; lin < ch.Line-1; lin++ {
				out.Print(" ",newer[lin])
			}
			out.Print(ch)
			lin += len(ch.New)
			for lin=lin; lin<len(newer)-1 && lin < ch.Line+len(ch.New)+3-1;lin++ {
				out.Print(" ",newer[lin])
			}
			c,e := out.PromptForChar(goopt.Expand("Verb this change? "))
			error.FailOn(e)
			switch c {
			case 'q','Q': error.Exit(e)
			case 'y','Y':
				debug.Println("Dealing with file ",filename)
				fstnum := ch.Line-1+offset
				start := outer[0:fstnum]
				end := outer[fstnum+len(ch.Old):]
				outer = ss.Cat(make([]string,0,len(start)+len(end)+len(ch.New)), start)
				outer = ss.Cat(outer, ch.New)
				outer = ss.Cat(outer, end)
				continue chunks
			case 'n','N': debug.Println("Ignoring changes to file ",filename)
				offset += len(ch.Old) - len(ch.New)
				continue chunks
			}
		}
	}
	return strings.Join(outer,"")
}
