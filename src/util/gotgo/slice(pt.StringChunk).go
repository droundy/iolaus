package slice龍ptøStringChunk龍ptøStringChunk

import pt "../patienceTypes"


// Here we have some utility slice routines

// Map1 provides an in-place map, meaning it modifies its input slice.
// If you still want that data, use the Map function.
func Map1(f func(pt.StringChunk) pt.StringChunk, slice []pt.StringChunk) {
	for i,v := range slice {
		slice[i] = f(v)
	}
}

// Map provides an out-of-place map, meaning it does not modify its
// input slice.  It therefore has the advantage that you can Map from
// one type of slice to another.
func Map(f func(pt.StringChunk) pt.StringChunk, slice []pt.StringChunk) []pt.StringChunk {
	out := make([]pt.StringChunk, len(slice))
	for i,v := range slice {
		out[i] = f(v)
	}
	return out
}

func Fold(f func(pt.StringChunk, pt.StringChunk) pt.StringChunk, x pt.StringChunk, slice []pt.StringChunk) pt.StringChunk {
  for _, v := range slice {
    x = f(x, v)
  }
  return x
}

// Filter returns a slice containing only those elements for which the
// predicate function returns true.
func Filter(f func(pt.StringChunk) bool, slice []pt.StringChunk) []pt.StringChunk {
	out := make ([]pt.StringChunk, 0, len(slice))
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
func Append(slice []pt.StringChunk, val pt.StringChunk) []pt.StringChunk {
	length := len(slice)
	if cap(slice) == length {
		// we need to expand
		newsl := make([]pt.StringChunk, length, 2*(length+1))
		for i,v := range slice {
			newsl[i] = v
		}
		slice = newsl
	}
	slice = slice[0:length+1]
	slice[length] = val
	return slice
}

func Repeat(val pt.StringChunk, n int) []pt.StringChunk {
	out := make([]pt.StringChunk, n)
	for i,_ := range out { out[i] = val }
	return out
}

// Cat concatenates two slices, expanding if needed.
func Cat(slice []pt.StringChunk, more []pt.StringChunk) []pt.StringChunk {
	length1, length2 := len(slice), len(more)
	if cap(slice) <= length1 + length2 {
		// we need to expand
		newsl := make([]pt.StringChunk, length1, 2*(length1+length2))
		for i,v := range slice {
			newsl[i] = v
		}
		slice = newsl
	}
	slice = slice[0:length1+length2]
	for i,v := range more {
		slice[length1+i] = v
	}
	return slice
}

func Reverse(slice []pt.StringChunk) (out []pt.StringChunk) {
	ln := len(slice)
	out = make([]pt.StringChunk, ln)
	for i,v:= range slice {
		out[ln-1-i] = v
	}
	return
}
// Here we will test that the types parameters are ok...


func testTypes(arg0 pt.StringChunk, arg1 pt.StringChunk) {
    f := func(interface{}, interface{}) { } // this func does nothing...
    f(arg0, arg1)
}
