package xgen

import (
	"github.com/goxgen/goxgen/projects"
	"regexp"
)

type XgenOption = func(xgen *Xgen) error

// WithProject adds a new project to Xgen
func WithProject(name string, project projects.Project) XgenOption {
	valid := regexp.MustCompile("^[a-z][a-z0-9_]*$").MatchString(name)
	if !valid {
		panic("project name must be a valid go identifier, \"%s\" provided")
	}
	return func(x *Xgen) (err error) {
		x.Projects[name] = project
		return nil
	}
}

// WithPackageName sets package name for Xgen
func WithPackageName(packageName string) XgenOption {
	return func(x *Xgen) (err error) {
		x.PackageName = &packageName
		return nil
	}
}
