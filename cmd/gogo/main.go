// Copyright 2020 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/ggaaooppeenngg/gogo/go/go2go"
)

var gotool = filepath.Join(runtime.GOROOT(), "bin", "go")

var cmds = map[string]bool{
	"build":     true,
	"run":       true,
	"test":      true,
	"translate": true,
}

// tagsFlag is the implementation of the -tags flag.
type tagsFlag []string

var buildTags tagsFlag

func (v *tagsFlag) Set(s string) error {
	// Split on commas, ignore empty strings.
	*v = []string{}
	for _, s := range strings.Split(s, ",") {
		if s != "" {
			*v = append(*v, s)
		}
	}
	return nil
}

func (v *tagsFlag) String() string {
	return strings.Join(*v, ",")
}

func main() {
	flag.Usage = usage
	flag.Parse()

	args := flag.Args()
	if len(args) < 1 {
		usage()
	}

	if !cmds[args[0]] {
		usage()
	}
	cmd := args[0]

	fs := flag.NewFlagSet("", flag.ContinueOnError)
	fs.Var((*tagsFlag)(&buildTags), "tags", "tag,list")
	fs.Parse(args[1:])

	args = fs.Args()

	importerTmpdir, err := ioutil.TempDir("", "go2go")
	if err != nil {
		log.Fatal(err)
	}
	defer os.RemoveAll(importerTmpdir)

	importer := go2go.NewImporter(importerTmpdir)

	if len(buildTags) > 0 {
		importer.SetTags(buildTags)
	}

	var rundir string
	if cmd == "run" {
		tmpdir := copyToTmpdir(args)
		defer os.RemoveAll(tmpdir)
		translate(importer, tmpdir)
		nargs := []string{"run"}
		for _, arg := range args {
			base := filepath.Base(arg)
			f := strings.TrimSuffix(base, ".go2") + ".go"
			nargs = append(nargs, f)
		}
		args = nargs
		rundir = tmpdir
	} else if cmd == "translate" && isGo2Files(args...) {
		for _, arg := range args {
			translateFile(importer, arg)
		}
	} else {
		for _, dir := range expandPackages(args) {
			translate(importer, dir)
		}
	}

	if cmd != "translate" {
		if len(buildTags) > 0 {
			args = append([]string{
				fmt.Sprintf("-tags=%s", strings.Join(buildTags, ",")),
			}, args...)
		}

		args = append([]string{cmd}, args...)

		cmd := exec.Command(gotool, args...)
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Dir = rundir
		gopath := importerTmpdir
		if go2path := os.Getenv("GO2PATH"); go2path != "" {
			gopath += string(os.PathListSeparator) + go2path
		}
		if oldGopath := os.Getenv("GOPATH"); oldGopath != "" {
			gopath += string(os.PathListSeparator) + oldGopath
		}
		cmd.Env = append(os.Environ(),
			"GOPATH="+gopath,
			"GO111MODULE=off",
		)
		if err := cmd.Run(); err != nil {
			die(fmt.Sprintf("%s %v failed: %v", gotool, args, err))
		}
	}
}

// isGo2Files reports whether the arguments are a list of .go2 files.
func isGo2Files(args ...string) bool {
	for _, arg := range args {
		if filepath.Ext(arg) != ".go2" {
			return false
		}
	}
	return true
}

// expandPackages returns a list of directories expanded from packages.
func expandPackages(pkgs []string) []string {
	if len(pkgs) == 0 {
		return []string{"."}
	}
	go2path := os.Getenv("GO2PATH")
	var dirs []string
pkgloop:
	for _, pkg := range pkgs {
		if go2path != "" {
			for _, pd := range strings.Split(go2path, string(os.PathListSeparator)) {
				d := filepath.Join(pd, "src", pkg)
				if fi, err := os.Stat(d); err == nil && fi.IsDir() {
					dirs = append(dirs, d)
					continue pkgloop
				}
			}
		}

		cmd := exec.Command(gotool, "list", "-f", "{{.Dir}}", pkg)
		cmd.Stderr = os.Stderr
		if go2path != "" {
			gopath := go2path
			if oldGopath := os.Getenv("GOPATH"); oldGopath != "" {
				gopath += string(os.PathListSeparator) + oldGopath
			}
			cmd.Env = append(os.Environ(),
				"GOPATH="+gopath,
				"GO111MODULE=off",
			)
		}
		out, err := cmd.Output()
		if err != nil {
			die(fmt.Sprintf("%s list %q failed: %v", gotool, pkg, err))
		}
		dirs = append(dirs, strings.Split(string(out), "\n")...)
	}
	return dirs
}

// copyToTmpdir copies files into a temporary directory.
func copyToTmpdir(files []string) string {
	if len(files) == 0 {
		die("no files to run")
	}
	tmpdir, err := ioutil.TempDir("", "go2go-run")
	if err != nil {
		die(err.Error())
	}
	for _, file := range files {
		data, err := ioutil.ReadFile(file)
		if err != nil {
			die(err.Error())
		}
		if err := ioutil.WriteFile(filepath.Join(tmpdir, filepath.Base(file)), data, 0444); err != nil {
			die(err.Error())
		}
	}
	return tmpdir
}

// usage reports a usage message and exits with failure.
func usage() {
	fmt.Fprint(os.Stderr, `Usage: go2go <command> [arguments]

The commands are:

	build      translate and build packages
	run        translate and run list of files
	test       translate and test packages
	translate  translate .go2 files into .go files
`)
	os.Exit(2)
}

// die reports an error and exits.
func die(msg string) {
	fmt.Fprintln(os.Stderr, msg)
	os.Exit(1)
}
