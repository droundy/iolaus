package slice

import a "../git/git"


// Here we have some utility slice routines

// Map1 provides an in-place map, meaning it modifies its input slice.
// If you still want that data, use the Map function.
func Map1(f func(a.Commitish) a.Commitish, slice []a.Commitish) {
	for i,v := range slice {
		slice[i] = f(v)
	}
}

// Map provides an out-of-place map, meaning it does not modify its
// input slice.  It therefore has the advantage that you can Map from
// one type of slice to another.
func Map(f func(a.Commitish) a.Commitish, slice []a.Commitish) []a.Commitish {
	out := make([]a.Commitish, len(slice))
	for i,v := range slice {
		out[i] = f(v)
	}
	return out
}

func Fold(f func(a.Commitish, a.Commitish) a.Commitish, x a.Commitish, slice []a.Commitish) a.Commitish {
  for _, v := range slice {
    x = f(x, v)
  }
  return x
}

// Filter returns a slice containing only those elements for which the
// predicate function returns true.
func Filter(f func(a.Commitish) bool, slice []a.Commitish) []a.Commitish {
	out := make ([]a.Commitish, 0, len(slice))
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
func Append(slice []a.Commitish, val a.Commitish) []a.Commitish {
	length := len(slice)
	if cap(slice) == length {
		// we need to expand
		newsl := make([]a.Commitish, length, 2*(length+1))
		for i,v := range slice {
			newsl[i] = v
		}
		slice = newsl
	}
	slice = slice[0:length+1]
	slice[length] = val
	return slice
}

func Repeat(val a.Commitish, n int) []a.Commitish {
	out := make([]a.Commitish, n)
	for i,_ := range out { out[i] = val }
	return out
}

// Cat concatenates two slices, expanding if needed.
func Cat(slices ...[]a.Commitish) []a.Commitish {
	return Cats(slices)
}

// Cats concatenates several slices, expanding if needed.
func Cats(slices [][]a.Commitish) []a.Commitish {
	lentot := 0
	for _,sl := range slices {
		lentot += len(sl)
	}
	out := make([]a.Commitish, lentot)
	i := 0
	for _,sl := range slices {
		for _,v := range sl {
			out[i] = v
			i++
		}
	}
	return out
}

func Reverse(slice []a.Commitish) (out []a.Commitish) {
	ln := len(slice)
	out = make([]a.Commitish, ln)
	for i,v:= range slice {
		out[ln-1-i] = v
	}
	return
}

func Any(f func(a.Commitish) bool, slice []a.Commitish) bool {
	for _,v:= range slice {
		if f(v) { return true }
	}
	return false
}

// Here we will test that the types parameters are ok...
func testTypes(arg0 a.Commitish, arg1 a.Commitish) {
    f := func(interface{}, interface{}) { } // this func does nothing...
    f(arg0, arg1)
}
