package goxgen

// ProjectOption is a function that configures a Xgen
type ProjectOption = func(project *SimpleProject) error

// WithOutputDir sets the output directory of generated code
func WithOutputDir(outputDir string) ProjectOption {
	return func(g *SimpleProject) (err error) {
		g.outputDir = StringP(outputDir)
		return nil
	}
}
