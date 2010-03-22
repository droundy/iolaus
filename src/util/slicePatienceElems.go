package patience



// Here we have some utility slice routines

// Map1 provides an in-place map, meaning it modifies its input slice.
// If you still want that data, use the Map function.
func pesMap1(f func([]PatienceElem) []PatienceElem, slice [][]PatienceElem) {
	for i,v := range slice {
		slice[i] = f(v)
	}
}

// Map provides an out-of-place map, meaning it does not modify its
// input slice.  It therefore has the advantage that you can Map from
// one type of slice to another.
func pesMap(f func([]PatienceElem) []PatienceElem, slice [][]PatienceElem) [][]PatienceElem {
	out := make([][]PatienceElem, len(slice))
	for i,v := range slice {
		out[i] = f(v)
	}
	return out
}

func pesFold(f func([]PatienceElem, []PatienceElem) []PatienceElem, x []PatienceElem, slice [][]PatienceElem) []PatienceElem {
  for _, v := range slice {
    x = f(x, v)
  }
  return x
}

// Filter returns a slice containing only those elements for which the
// predicate function returns true.
func pesFilter(f func([]PatienceElem) bool, slice [][]PatienceElem) [][]PatienceElem {
	out := make ([][]PatienceElem, 0, len(slice))
	i := 0
	for _,v := range slice {
		if f(v) {
			out = out[0:i+1]
			out[i] = v
			i++
		}
	}
	return out
}

// Append appends an element to a slice, in-place if possible, and
// expanding if needed.
func pesAppend(slice [][]PatienceElem, val []PatienceElem) [][]PatienceElem {
	length := len(slice)
	if cap(slice) == length {
		// we need to expand
		newsl := make([][]PatienceElem, length, 2*(length+1))
		for i,v := range slice {
			newsl[i] = v
		}
		slice = newsl
	}
	slice = slice[0:length+1]
	slice[length] = val
	return slice
}

func pesRepeat(val []PatienceElem, n int) [][]PatienceElem {
	out := make([][]PatienceElem, n)
	for i,_ := range out { out[i] = val }
	return out
}

// Cat concatenates two slices, expanding if needed.
func pesCat(slices ...[][]PatienceElem) [][]PatienceElem {
	return pesCats(slices)
}

// Cats concatenates several slices, expanding if needed.
func pesCats(slices [][][]PatienceElem) [][]PatienceElem {
	lentot := 0
	for _,sl := range slices {
		lentot += len(sl)
	}
	out := make([][]PatienceElem, lentot)
	i := 0
	for _,sl := range slices {
		for _,v := range sl {
			out[i] = v
			i++
		}
	}
	return out
}

func pesReverse(slice [][]PatienceElem) (out [][]PatienceElem) {
	ln := len(slice)
	out = make([][]PatienceElem, ln)
	for i,v:= range slice {
		out[ln-1-i] = v
	}
	return
}

func pesAny(f func([]PatienceElem) bool, slice [][]PatienceElem) bool {
	for _,v:= range slice {
		if f(v) { return true }
	}
	return false
}

// Here we will test that the types parameters are ok...
func pestestTypes(arg0 []PatienceElem, arg1 []PatienceElem) {
    f := func(interface{}, interface{}) { } // this func does nothing...
    f(arg0, arg1)
}
