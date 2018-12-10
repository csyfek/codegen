package main

import (
	"github.com/jackmanlabs/errors"
	"go/build"
	"log"
	"os"
	"path"
	"strings"
)

func findPackagePaths() ([]*PackageTreeItem, error) {
	gopath := build.Default.GOPATH

	root, err := findPackagePathsRecurse(gopath + "/src")
	if err != nil {
		return nil, errors.Stack(err)
	}

	root._packageName = "GOPATH"

	return []*PackageTreeItem{root}, nil
}

func findPackagePathsRecurse(dir string) (*PackageTreeItem, error) {

	if path.Base(dir) == "internal" {
		return nil, nil
	}

	if strings.HasPrefix(path.Base(dir), ".") {
		return nil, nil
	}

	f, err := os.Open(dir)
	if err != nil {
		return nil, errors.Stack(err)
	}

	children := make([]*PackageTreeItem, 0)

	subs, err := f.Readdir(-1)
	for _, sub := range subs {

		if sub.IsDir() {
			newChild, err := findPackagePathsRecurse(dir + "/" + sub.Name())
			if err != nil {
				return nil, errors.Stack(err)
			}

			if newChild != nil {
				children = append(children, newChild)
			}
		}
	}

	backup := NewPackageTreeItem(nil).initWith(path.Base(dir), "")

	p, err := build.ImportDir(dir, 0)
	if err != nil {
		switch err.(type) {
		case *build.NoGoError:
		// Do nothing. This error is acceptable.
		default:
			log.Print(err)
		}
		if len(children) == 0 {
			return nil, nil
		}
		return backup.withChildren(children), nil
	}

	if p.Name == "main" {
		if len(children) == 0 {
			return nil, nil
		}
		return backup.withChildren(children), nil
	}

	i := NewPackageTreeItem(nil).initWith(p.Name, p.ImportPath).withChildren(children)

	return i, nil
}
