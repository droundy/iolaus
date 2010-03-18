package slice龍øøptøPatienceElem龍øøptøPatienceElem

import pt "../patienceTypes"


// Here we have some utility slice routines

// Map1 provides an in-place map, meaning it modifies its input slice.
// If you still want that data, use the Map function.
func Map1(f func([]pt.PatienceElem) []pt.PatienceElem, slice [][]pt.PatienceElem) {
	for i,v := range slice {
		slice[i] = f(v)
	}
}

// Map provides an out-of-place map, meaning it does not modify its
// input slice.  It therefore has the advantage that you can Map from
// one type of slice to another.
func Map(f func([]pt.PatienceElem) []pt.PatienceElem, slice [][]pt.PatienceElem) [][]pt.PatienceElem {
	out := make([][]pt.PatienceElem, len(slice))
	for i,v := range slice {
		out[i] = f(v)
	}
	return out
}

func Fold(f func([]pt.PatienceElem, []pt.PatienceElem) []pt.PatienceElem, x []pt.PatienceElem, slice [][]pt.PatienceElem) []pt.PatienceElem {
  for _, v := range slice {
    x = f(x, v)
  }
  return x
}

// Filter returns a slice containing only those elements for which the
// predicate function returns true.
func Filter(f func([]pt.PatienceElem) bool, slice [][]pt.PatienceElem) [][]pt.PatienceElem {
	out := make ([][]pt.PatienceElem, 0, len(slice))
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
func Append(slice [][]pt.PatienceElem, val []pt.PatienceElem) [][]pt.PatienceElem {
	length := len(slice)
	if cap(slice) == length {
		// we need to expand
		newsl := make([][]pt.PatienceElem, length, 2*(length+1))
		for i,v := range slice {
			newsl[i] = v
		}
		slice = newsl
	}
	slice = slice[0:length+1]
	slice[length] = val
	return slice
}

func Repeat(val []pt.PatienceElem, n int) [][]pt.PatienceElem {
	out := make([][]pt.PatienceElem, n)
	for i,_ := range out { out[i] = val }
	return out
}

// Cat concatenates two slices, expanding if needed.
func Cat(slices ...[][]pt.PatienceElem) [][]pt.PatienceElem {
	return Cats(slices)
}

// Cats concatenates several slices, expanding if needed.
func Cats(slices [][][]pt.PatienceElem) [][]pt.PatienceElem {
	lentot := 0
	for _,sl := range slices {
		lentot += len(sl)
	}
	out := make([][]pt.PatienceElem, lentot)
	i := 0
	for _,sl := range slices {
		for _,v := range sl {
			out[i] = v
			i++
		}
	}
	return out
}

func Reverse(slice [][]pt.PatienceElem) (out [][]pt.PatienceElem) {
	ln := len(slice)
	out = make([][]pt.PatienceElem, ln)
	for i,v:= range slice {
		out[ln-1-i] = v
	}
	return
}

func Any(f func([]pt.PatienceElem) bool, slice [][]pt.PatienceElem) bool {
	for _,v:= range slice {
		if f(v) { return true }
	}
	return false
}
// Here we will test that the types parameters are ok...


func testTypes(arg0 []pt.PatienceElem, arg1 []pt.PatienceElem) {
    f := func(interface{}, interface{}) { } // this func does nothing...
    f(arg0, arg1)
}
