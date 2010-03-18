package slice龍gitøCommitHash龍gitøCommitHash

import git "../git"


// Here we have some utility slice routines

// Map1 provides an in-place map, meaning it modifies its input slice.
// If you still want that data, use the Map function.
func Map1(f func(git.CommitHash) git.CommitHash, slice []git.CommitHash) {
	for i,v := range slice {
		slice[i] = f(v)
	}
}

// Map provides an out-of-place map, meaning it does not modify its
// input slice.  It therefore has the advantage that you can Map from
// one type of slice to another.
func Map(f func(git.CommitHash) git.CommitHash, slice []git.CommitHash) []git.CommitHash {
	out := make([]git.CommitHash, len(slice))
	for i,v := range slice {
		out[i] = f(v)
	}
	return out
}

func Fold(f func(git.CommitHash, git.CommitHash) git.CommitHash, x git.CommitHash, slice []git.CommitHash) git.CommitHash {
  for _, v := range slice {
    x = f(x, v)
  }
  return x
}

// Filter returns a slice containing only those elements for which the
// predicate function returns true.
func Filter(f func(git.CommitHash) bool, slice []git.CommitHash) []git.CommitHash {
	out := make ([]git.CommitHash, 0, len(slice))
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
func Append(slice []git.CommitHash, val git.CommitHash) []git.CommitHash {
	length := len(slice)
	if cap(slice) == length {
		// we need to expand
		newsl := make([]git.CommitHash, length, 2*(length+1))
		for i,v := range slice {
			newsl[i] = v
		}
		slice = newsl
	}
	slice = slice[0:length+1]
	slice[length] = val
	return slice
}

func Repeat(val git.CommitHash, n int) []git.CommitHash {
	out := make([]git.CommitHash, n)
	for i,_ := range out { out[i] = val }
	return out
}

// Cat concatenates two slices, expanding if needed.
func Cat(slices ...[]git.CommitHash) []git.CommitHash {
	return Cats(slices)
}

// Cats concatenates several slices, expanding if needed.
func Cats(slices [][]git.CommitHash) []git.CommitHash {
	lentot := 0
	for _,sl := range slices {
		lentot += len(sl)
	}
	out := make([]git.CommitHash, lentot)
	i := 0
	for _,sl := range slices {
		for _,v := range sl {
			out[i] = v
			i++
		}
	}
	return out
}

func Reverse(slice []git.CommitHash) (out []git.CommitHash) {
	ln := len(slice)
	out = make([]git.CommitHash, ln)
	for i,v:= range slice {
		out[ln-1-i] = v
	}
	return
}

func Any(f func(git.CommitHash) bool, slice []git.CommitHash) bool {
	for _,v:= range slice {
		if f(v) { return true }
	}
	return false
}
// Here we will test that the types parameters are ok...


func testTypes(arg0 git.CommitHash, arg1 git.CommitHash) {
    f := func(interface{}, interface{}) { } // this func does nothing...
    f(arg0, arg1)
}
