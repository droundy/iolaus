// Copyright 2009 Dimiter Stanev, malkia@gmail.com.
// Copyright 2010 David Roundy, roundyd@physics.oregonstate.edu.
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"os"
	"fmt"
	"io/ioutil"
	"go/parser"
	"go/ast"
	"strconv"
	"path"
        "sort"
)

var (
	curdir, _ = os.Getwd()
)

func getImports(filename string) (pkgname string,
	imports map[string]bool, names map[string]string, error os.Error) {
	source, error := ioutil.ReadFile(filename)
	if error != nil {
		return
	}
	file, error := parser.ParseFile(filename, source, nil, parser.ImportsOnly)
	if error != nil {
		return
	}
	dir, _ := path.Split(filename)
	pkgname = file.Name.Name()
	for _, importDecl := range file.Decls {
		importDecl, ok := importDecl.(*ast.GenDecl)
		if ok {
			for _, importSpec := range importDecl.Specs {
				importSpec, ok := importSpec.(*ast.ImportSpec)
				if ok {
					importPath, _ := strconv.Unquote(string(importSpec.Path.Value))
					if len(importPath) > 0 {
						if imports == nil {
							imports = make(map[string]bool)
							names = make(map[string]string)
						}
						if importPath[0] == '.' {
							imports[path.Join(dir, path.Clean(importPath))] = true
						}
					}
				}
			}
		}
	}
	return
}

func cleanbinname(f string) string {
	if f[0:4] == "src/" { return "bin/"+f[4:] }
	return f
}

func shouldUpdate(sourceFile, targetFile string) (doUpdate bool, error os.Error) {
	sourceStat, error := os.Lstat(sourceFile)
	if error != nil {
		return false, error
	}
	targetStat, error := os.Lstat(targetFile)
	if error != nil {
		return true, error
	}
	return targetStat.Mtime_ns < sourceStat.Mtime_ns, error
}

type maker struct {
}
func (maker) VisitDir(string, *os.Dir) bool { return true }
func (maker) VisitFile(f string, _ *os.Dir) {
	if path.Ext(f) == ".go" {
		deps := make([]string, 1, 1000) // FIXME stupid hardcoded limit...x
		deps[0] = f
		pname, imports, _, err := getImports(f)
		if err != nil {
			fmt.Println("# error: ", err)
		}
		for i,_ := range imports {
			deps = deps[0:len(deps)+1]
			deps[len(deps)-1] = i + ".$(O)"
		}
                sort.SortStrings(deps[1:]) // alphebatize all deps but first
		basename := f[0:len(f)-3]
		objname := basename+".$(O)"
		if pname == "main" {
			fmt.Printf("%s: %s\n\t@mkdir -p bin\n\t$(LD) -o $@ $<\n", cleanbinname(basename), objname)
		}
		fmt.Print(objname+":")
		for _,d := range deps {
			fmt.Print(" "+d)
		}
		fmt.Print("\n\n")
	}
}

type seeker struct {
}
func (seeker) VisitDir(string, *os.Dir) bool { return true }
func (seeker) VisitFile(f string, _ *os.Dir) {
	if path.Ext(f) == ".go" {
		pname, _, _, _ := getImports(f)
		basename := f[0:len(f)-3]
		if pname == "main" {
			mybinfiles += " " + cleanbinname(basename)
		}
	}
}

var mybinfiles = ""

func main() {
	fmt.Print(`# Copyright 2010 David Roundy, roundyd@physics.oregonstate.edu.
# All rights reserved.

`)
	path.Walk(".", seeker{}, nil)
	fmt.Printf(`
all: Makefile %s

include $(GOROOT)/src/Make.$(GOARCH)

.PHONY: test
.SUFFIXES: .$(O) .go .got .gotgo

`, mybinfiles)
	fmt.Print("Makefile: scripts/mkmake\n\t./scripts/mkmake > $@\n")
	fmt.Print(".go.$(O):\n\tcd `dirname \"$<\"`; $(GC) `basename \"$<\"`\n")
	fmt.Print(".got.gotgo:\n\tgotgo \"$<\"\n\n")
	path.Walk(".", maker{}, nil)
}
