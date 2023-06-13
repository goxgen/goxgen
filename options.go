package goxgen

type XgenOption = func(xgen *Xgen) error

// WithProject adds a new project to Xgen
func WithProject(name string, option ...ProjectOption) XgenOption {
	return func(x *Xgen) (err error) {
		x.Projects = append(x.Projects, NewSimpleProject(name, option...))
		return nil
	}
}

// WithEntProject adds a new Ent project to Xgen
func WithEntProject(name string, option ...ProjectOption) XgenOption {
	return func(x *Xgen) (err error) {
		x.Projects = append(x.Projects, NewEntProject(name, option...))
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
