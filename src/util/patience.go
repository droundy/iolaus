package patience

import (
	pt "./patienceTypes"
	intslice "./gotgo/slice(int)"
	ch "./gotgo/slice(pt.StringChunk)"
	pes "./gotgo/slice(pt.PatienceElem)"
	pess "./gotgo/slice([]pt.PatienceElem)"
)

func Diff(o, n []string) []pt.StringChunk {
	return DiffFromLine(1,o,n)
}

func DiffFromLine(line0 int, o, n []string) []pt.StringChunk {
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
	uniques := []int{}
	for onum, l := range o {
		if nnum, ok := nnums[l]; ok && nnum != 0 {
			uniques = intslice.Append(uniques, nnum)
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
			ostuff := []string{}
			if len(o)-last > first {
				ostuff = o[first:len(o)-last]
			}
			nstuff := []string{}
			if len(n)-last > first {
				nstuff = n[first:len(n)-last]
			}
			return []pt.StringChunk{pt.StringChunk{line0+first, ostuff, nstuff}}
		} else {
			return []pt.StringChunk{}
		}
	}
	piles := [][]pt.PatienceElem{[]pt.PatienceElem{pt.PatienceElem{uniques[0],0}}}
	for _,v := range uniques[1:] {
		foundone := false
		for ipile,pile := range piles {
			if pile[0].Val > v {
				foundone = true
				var myprev int
				if ipile > 0 {
					// points to top element of previous pile 
					myprev = len(piles[ipile-1])-1
				}
				piles[ipile] = pes.Append(piles[ipile], pt.PatienceElem{v,myprev})
				break
			}
		}
		if !foundone {
			newpile := make([]pt.PatienceElem,1,4)
			newpile[0] = pt.PatienceElem{v, len(piles[len(piles)-1])-1}
			//fmt.Println("len(piles) is",len(piles))
			//fmt.Println("cap(piles) is",cap(piles))
			piles = pess.Append(piles, newpile)
		}
	}
	//fmt.Println("unique are",uniques)
	//fmt.Println("piles are",piles)
	lcs := []int{}
	for pnum, enum := len(piles)-1, len(piles[len(piles)-1])-1; pnum >= 0; pnum-- {
		//fmt.Println("pnum is",pnum)
		//fmt.Println("enum is",enum)
		//fmt.Println("len(piles[pnum])",len(piles[pnum]))
		lcs = intslice.Append(lcs, piles[pnum][enum].Val)
		enum = piles[pnum][enum].Prev
	}
	diff := []pt.StringChunk{}
	for prevo,prevn,i:= 0,0,len(lcs)-1; i>=0; i-- {
		nextn := lcs[i]
		//fmt.Println("looking for",n[nextn])
		nexto := onums[n[nextn]]
		//fmt.Println("nexto",nexto)
		//fmt.Println("uniques is",uniques)
		//fmt.Println("nnums is",nnums)
		diff = ch.Cat(diff,
			DiffFromLine(line0+prevn, o[prevo:nexto], n[prevn:nextn]))
		prevo = nexto
		prevn = nextn
	}
	return diff
}
