package patienceTypes

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
		out += "-" + l
	}
	for _, l := range ch.New {
		out += "+" + l
	}
	return
}

type PatienceElem struct {
	Val int
	// The prev points to the next "pile" to read off.
	Prev int
}
