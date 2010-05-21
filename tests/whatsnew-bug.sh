#!/bin/sh

set -ev

# This is a test of whatsnew...

mkdir test1
cd test1
iolaus-initialize

cat > file <<EOF
		e = plumbing.ReadTree(local, "--index-output=.git/index.pulling")
		error.FailOn(e)
		e = os.Setenv("GIT_INDEX_FILE", ".git/index.pulling")
		error.FailOn(e)
		t, e := core.Merge(local, remote)
		error.FailOn(e)
		c := plumbing.CommitTree(t, []git.Commitish{local,remote}, "Foo")
		plumbing.UpdateRef("HEAD", c)
		plumbing.CheckoutIndex("--all")
		// Now let's update the true index by just copying over the scratch...
		e = os.Rename(".git/index.pulling", ".git/index")
		error.FailOn(e) // FIXME: we should do better than this...
		exit.Exit(0)
EOF

iolaus-record -am 'add file'

cat > file <<EOF
		e = plumbing.ReadTree(local, "--index-output=.git/index.pulling")
		error.FailOn(e)
		e = os.Setenv("GIT_INDEX_FILE", ".git/index.pulling")
		error.FailOn(e)
		debug.Println("About to merge", local, "with", remote)
		t, e := core.Merge(local, remote)
		debug.Println("Done merging.")
		error.FailOn(e)
		c := plumbing.CommitTree(t, []git.Commitish{local,remote}, "Foo")
		plumbing.UpdateRef("HEAD", c)
		plumbing.CheckoutIndex("--all")
		// Now let's update the true index by just copying over the scratch...
		e = os.Rename(".git/index.pulling", ".git/index")
		error.FailOn(e) // FIXME: we should do better than this...
		exit.Exit(0)
EOF

iolaus-whatsnew

iolaus-whatsnew > out

# To be careful, I'll look at the debug message "Done merging."
grep -1 'Done merging' out

# This line should only show up once:
grep 'Done merging' out | wc -l | grep 1

# And it should be preceded by a Merge:
grep -1 'Done merging' out | grep Merge

# And it should be followed by a FailOn:
grep -1 'Done merging' out | grep FailOn


# This particular diff showed problems around the added line that
# added the debug message: "About to merge".
grep -1 About out

# About to merge should be followed by Merge
grep -1 About out | grep Merge

# About to merge should be preceded by FailOn
grep -1 About out | grep FailOn

# "About to merge" should occur only once
grep About out | wc -l | grep 1
