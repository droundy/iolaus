package patience

import (
	//"strconv"
	//"fmt"
	"io/ioutil"
	"strings"
	"./debug"
	pt "./patienceTypes"
	intslice "./gotgo/slice(int)"
	ch "./gotgo/slice(pt.StringChunk)"
	pes "./gotgo/slice(pt.PatienceElem)"
	pess "./gotgo/slice([]pt.PatienceElem)"
)

func Diff(o, n []string) []pt.StringChunk {
	ioutil.WriteFile("/tmp/old", []byte(strings.Join(o,"")), 0666)
	ioutil.WriteFile("/tmp/new", []byte(strings.Join(n,"")), 0666)
	return DiffFromLine(1,o,n)
}

func DiffFromLine(line0 int, o, n []string) []pt.StringChunk {
	if len(o) == 0 && len(n) == 0 {
		return []pt.StringChunk{}
	}
	nnums := map[string] int{}
	for lnum, l := range n {
		_, present := nnums[l]
		if present {
			nnums[l] = -1 // use -1 to indicate it's a duplicate
		} else {
			nnums[l] = lnum // this is distance from "line0"
		}
		//debug.Print("n: ", l)
	}
	onums := map[string]int {}
	for lnum, l := range o {
		_, present := onums[l]
		if present {
			onums[l] = -1 // use -1 to indicate it's a duplicate
		} else {
			onums[l] = lnum // this is distance from "line0"
		}
		//debug.Print("o: ", l)
	}
	uniques := []int{}
	for _, l := range o {
		nnum, okn := nnums[l]
		onum, oko := onums[l]
		if oko && okn && nnum != -1 && onum != -1 {
			uniques = intslice.Append(uniques, nnum)
		}
	}
	debug.Printf("Uniques are %v\n", uniques)
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
			//fmt.Printf("Hello silly %v %v\n", len(o)-last > first, len(n)-last > first)
			//fmt.Printf("Hello lens %d %d %d %d\n", len(o),len(n),last, first)
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
			piles = pess.Append(piles, newpile)
		}
	}
	debug.Printf("Piles are %v\n", piles)
	lcs := []int{}
	for pnum, enum := len(piles)-1, len(piles[len(piles)-1])-1; pnum >= 0; pnum-- {
		lcs = intslice.Append(lcs, piles[pnum][enum].Val)
		enum = piles[pnum][enum].Prev
	}
	diff := []pt.StringChunk{}
	for prevo,prevn,i:= 0,0,len(lcs)-1; i>=0; i-- {
		nextn := lcs[i]
		nexto := onums[n[nextn]]
		//fmt.Printf("len(o)=%d len(n)=%d prevo=%d nexto=%d prevn=%d nextn=%d\n",
		//	len(o), len(n), prevo, nexto, prevn, nextn)
		debug.Printf("Looking at changes in old from %d to %d\n", prevo, nexto)
		debug.Printf("Looking at changes in new from %d to %d\n", prevn, nextn)
		diff = ch.Cat(diff,
			DiffFromLine(line0+prevn, o[prevo:nexto], n[prevn:nextn]))
		prevo = nexto+1
		prevn = nextn+1
	}
	lastn := lcs[0]
	lasto := onums[n[lastn]]
	diff = ch.Cat(diff, DiffFromLine(line0+lastn, o[lasto+1:], n[lastn+1:]))
	return diff
}
