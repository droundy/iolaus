package patienceTypes

import "../git/color"

type StringChunk struct {
    Line int
    Old  []string
    New  []string
}
func (ch StringChunk) String() (out string) {
	if len(ch.Old) == 0 && len(ch.New) == 0 {
		return ""
	}
	out = "" // fmt.Sprintln(" ",ch.Line)
	for _, l := range ch.Old {
		out += color.String("-" + l, color.Old)
	}
	for _, l := range ch.New {
		out += color.String("+" + l, color.New)
	}
	return
}

type PatienceElem struct {
	Val int
	// The prev points to the next "pile" to read off.
	Prev int
}
