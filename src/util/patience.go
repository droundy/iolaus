package patience

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

func Diff(o, n []string) []StringChunk {
	return DiffFromLine(1,o,n)
}

func DiffFromLine(line0 int, o, n []string) []StringChunk {
	nnums := map[string] int{}
	for lnum, l := range n {
		_, present := nnums[l]
		if present {
			nnums[l] = 0 // use 0 to indicate it's a duplicate
		} else {
			nnums[l] = lnum // this is distance from "line0"
		}
	}
	// construct a list of the "new" positions of unique elements of
	// the "old" array.
	onums := map[string] int{}
	uniques := make([]int,0,len(o)) // better to use min(len(o),len(n))...
	for onum, l := range o {
		if nnum, ok := nnums[l]; ok && nnum != 0 {
			uniques = uniques[0:len(uniques)+1]
			uniques[len(uniques)-1] = nnum
			onums[l] = onum // this is distance from "line0"
		}
	}
	if len(uniques) == 0 {
		first, last := 0, 0
		for first < len(o) && first < len(n) && o[first] == n[first] { first++ }
		for first < len(o) && first < len(n) &&
			len(o)-1-last >= 0 && len(n)-1-last >= 0 &&
			o[len(o)-1-last] == n[len(n)-1-last] { last++ }
		if len(o)-last > first || len(n)-last > first {
			return []StringChunk{StringChunk{line0+first,
					o[first:len(o)-last], n[first:len(n)-last]}}
		} else {
			return []StringChunk{}
		}
	}
	piles := make([][]patienceElem,1,len(uniques))
	piles[0] = []patienceElem{ patienceElem{uniques[0],0} }
	for _,v := range uniques[1:] {
		foundone := false
		for ipile,pile := range piles {
			if pile[0].val > v {
				foundone = true
				var myprev int
				if ipile > 0 {
					// points to top element of previous pile 
					myprev = len(piles[ipile-1])-1
				}
				append(patienceElem{v,myprev},&piles[ipile])
				break
			}
		}
		if !foundone {
			newpile := make([]patienceElem,1,4)
			newpile[0] = patienceElem{v, len(piles[len(piles)-1])-1}
			//fmt.Println("len(piles) is",len(piles))
			//fmt.Println("cap(piles) is",cap(piles))
			piles = piles[0:len(piles)+1]
			piles[len(piles)-1] = newpile
		}
	}
	//fmt.Println("unique are",uniques)
	//fmt.Println("piles are",piles)
	lcs := make([]int, 0, len(uniques))
	for pnum, enum := len(piles)-1, len(piles[len(piles)-1])-1; pnum >= 0; pnum-- {
		length := len(lcs)
		lcs = lcs[0:length+1]
		//fmt.Println("pnum is",pnum)
		//fmt.Println("enum is",enum)
		//fmt.Println("len(piles[pnum])",len(piles[pnum]))
		lcs[length] = piles[pnum][enum].val
		enum = piles[pnum][enum].prev
	}
	diff := make([]StringChunk, 0, 2*len(lcs))
	for prevo,prevn,i:= 0,0,len(lcs)-1; i>=0; i-- {
		nextn := lcs[i]
		//fmt.Println("looking for",n[nextn])
		nexto := onums[n[nextn]]
		//fmt.Println("nexto",nexto)
		//fmt.Println("uniques is",uniques)
		//fmt.Println("nnums is",nnums)
		join(&diff,
			DiffFromLine(line0+prevn, o[prevo:nexto], n[prevn:nextn]))
		prevo = nexto
		prevn = nextn
	}
	return diff
}

type patienceElem struct {
	val int
	// The prev points to the next "pile" to read off.
	prev int
}

func append(x patienceElem, slice *[]patienceElem) {
	length := len(*slice)
  if length + 1 > cap(*slice) {  // reallocate
    // Allocate double what's needed, for future growth.
    newSlice := make([]patienceElem, (length + 1)*2)
    // Copy data (could use bytes.Copy()).
    for i, c := range *slice {
      newSlice[i] = c
    }
    *slice = newSlice
  }
	*slice = (*slice)[0:length+1]
	(*slice)[length] = x
}

func join(big *[]StringChunk, little []StringChunk) {
	length := len(*big)
	llittle := len(little)
  if length + llittle > cap(*big) {  // reallocate
    // Allocate double what's needed, for future growth.
    newSlice := make([]StringChunk, length+llittle, (length + llittle)*2)
    // Copy data (could use bytes.Copy()).
    for i, c := range *big {
      newSlice[i] = c
    }
		//fmt.Println("newSlice is",newSlice)
		//fmt.Println("cap(newSlice) is",cap(newSlice))
    *big = newSlice
  }
	*big = (*big)[0:length+llittle]
	//fmt.Println("big is",*big)
	//fmt.Println("cap(big) is",cap(*big))
	for i,c := range little {
		(*big)[length+i] = c
	}
}
