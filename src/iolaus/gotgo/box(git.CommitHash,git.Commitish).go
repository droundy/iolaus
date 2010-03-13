package box龍gitøCommitHash龍gitøCommitish龍string

import git "../../git/git"


// Here we have some utility routines for handling "boxing"
// conversions from a to b, where a is assignment-compatible with b.

func Box(as []git.CommitHash) []git.Commitish {
	bs := make([]git.Commitish, len(as))
	for i,a := range as {
		bs[i] = a
	}
	return bs
}

func BoxMap(ma map[string]git.CommitHash) (mb map[string]git.Commitish) {
	for i,a := range ma {
		mb[i] = a
	}
	return
}
// Here we will test that the types parameters are ok...


func testTypes(arg0 git.CommitHash, arg1 git.Commitish, arg2 string) {
    f := func(interface{}, interface{}, string) { } // this func does nothing...
    f(arg0, arg1, string(arg2))
}
