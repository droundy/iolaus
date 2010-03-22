package box

import ta "../../git/git"
import tb "../../git/git"


// Here we have some utility routines for handling "boxing"
// conversions from a to b, where a is assignment-compatible with b.

func Box(as []ta.CommitHash) []tb.Commitish {
	bs := make([]tb.Commitish, len(as))
	for i,a := range as {
		bs[i] = a
	}
	return bs
}

func BoxMap(ma map[string]ta.CommitHash) (mb map[string]tb.Commitish) {
	for i,a := range ma {
		mb[i] = a
	}
	return
}

// Here we will test that the types parameters are ok...
func testTypes(arg0 ta.CommitHash, arg1 tb.Commitish, arg2 string) {
    f := func(interface{}, interface{}, string) { } // this func does nothing...
    f(arg0, arg1, string(arg2))
}
