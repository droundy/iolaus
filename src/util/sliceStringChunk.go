package patience



// Here we have some utility slice routines

// Map1 provides an in-place map, meaning it modifies its input slice.
// If you still want that data, use the Map function.
func scMap1(f func(StringChunk) StringChunk, slice []StringChunk) {
	for i,v := range slice {
		slice[i] = f(v)
	}
}

// Map provides an out-of-place map, meaning it does not modify its
// input slice.  It therefore has the advantage that you can Map from
// one type of slice to another.
func scMap(f func(StringChunk) StringChunk, slice []StringChunk) []StringChunk {
	out := make([]StringChunk, len(slice))
	for i,v := range slice {
		out[i] = f(v)
	}
	return out
}

func scFold(f func(StringChunk, StringChunk) StringChunk, x StringChunk, slice []StringChunk) StringChunk {
  for _, v := range slice {
    x = f(x, v)
  }
  return x
}

// Filter returns a slice containing only those elements for which the
// predicate function returns true.
func scFilter(f func(StringChunk) bool, slice []StringChunk) []StringChunk {
	out := make ([]StringChunk, 0, len(slice))
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
func scAppend(slice []StringChunk, val StringChunk) []StringChunk {
	length := len(slice)
	if cap(slice) == length {
		// we need to expand
		newsl := make([]StringChunk, length, 2*(length+1))
		for i,v := range slice {
			newsl[i] = v
		}
		slice = newsl
	}
	slice = slice[0:length+1]
	slice[length] = val
	return slice
}

func scRepeat(val StringChunk, n int) []StringChunk {
	out := make([]StringChunk, n)
	for i,_ := range out { out[i] = val }
	return out
}

// Cat concatenates two slices, expanding if needed.
func scCat(slices ...[]StringChunk) []StringChunk {
	return scCats(slices)
}

// Cats concatenates several slices, expanding if needed.
func scCats(slices [][]StringChunk) []StringChunk {
	lentot := 0
	for _,sl := range slices {
		lentot += len(sl)
	}
	out := make([]StringChunk, lentot)
	i := 0
	for _,sl := range slices {
		for _,v := range sl {
			out[i] = v
			i++
		}
	}
	return out
}

func scReverse(slice []StringChunk) (out []StringChunk) {
	ln := len(slice)
	out = make([]StringChunk, ln)
	for i,v:= range slice {
		out[ln-1-i] = v
	}
	return
}

func scAny(f func(StringChunk) bool, slice []StringChunk) bool {
	for _,v:= range slice {
		if f(v) { return true }
	}
	return false
}

// Here we will test that the types parameters are ok...
func sctestTypes(arg0 StringChunk, arg1 StringChunk) {
    f := func(interface{}, interface{}) { } // this func does nothing...
    f(arg0, arg1)
}
