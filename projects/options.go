package projects

// ProjectOption is a function that configures a Xgen
type ProjectOption = func(project Project) error
