package goxgen

type XgenOption = func(xgen *Xgen) error

func WithProject(name string, option ...ProjectOption) XgenOption {
	return func(x *Xgen) (err error) {
		x.Projects = append(x.Projects, NewProject(name, option...))
		return nil
	}
}

func WithEntProject(name string, option ...ProjectOption) XgenOption {
	return func(x *Xgen) (err error) {
		x.Projects = append(x.Projects, NewEntProject(name, option...))
		return nil
	}
}

func WithPackageName(packageName string) XgenOption {
	return func(x *Xgen) (err error) {
		x.PackageName = &packageName
		return nil
	}
}
