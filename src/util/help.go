package help

import (
	"github.com/droundy/goopt"
)

func Init(summary string, description func() string, getopts func() []string ) {
	goopt.Author = "David Roundy"
	goopt.Description = description
	goopt.Summary = summary
	goopt.Version = "0.0"
	goopt.Suite = "Iolaus"
	goopt.Parse(getopts)
}
