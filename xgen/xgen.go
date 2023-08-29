package xgen

import (
	"context"
	"fmt"
	"github.com/goxgen/goxgen/cli"
	"github.com/goxgen/goxgen/projects"
	"github.com/goxgen/goxgen/utils"
)

const (
	GeneratedFilePrefix = "generated_xgen_"
)

// Xgen main struct
type Xgen struct {
	PackageName *string                     `yaml:"package_name" json:"package_name"`
	Projects    map[string]projects.Project `yaml:"projects" json:"projects"`
	CLI         *cli.CLI                    `yaml:"cli" json:"cli"`
}

// NewXgen creates a new Xgen instance
// it creates a new Xgen instance with projects and CLI
func NewXgen(options ...XgenOption) *Xgen {
	xgen := &Xgen{
		CLI:      &cli.CLI{},
		Projects: map[string]projects.Project{},
	}

	for _, opt := range options {
		if err := opt(xgen); err != nil {
			panic(err)
		}
	}

	return xgen
}

// Generate generates code
// it is a main function of Xgen
// it generates code for projects and CLI
func (x *Xgen) Generate(ctx context.Context) (err error) {
	projCtx := projects.PrepareCommonContext(
		ctx,
		&projects.Context{
			ParentPackageName:   utils.PString(x.PackageName),
			GeneratedFilePrefix: GeneratedFilePrefix,
			Projects:            x.Projects,
		},
	)

	err = projects.NewProjectGenerator(x.Projects).Generate(projCtx)
	if err != nil {
		return fmt.Errorf("failed to generate projects: %w", err)
	}

	cliContext := cli.PrepareContext(
		ctx,
		&cli.Context{
			ParentPackageName:   utils.PString(x.PackageName),
			GeneratedFilePrefix: GeneratedFilePrefix,
			Projects:            x.Projects,
			OutputDir:           ".",
		},
	)
	err = x.CLI.Generate(cliContext)
	if err != nil {
		return err
	}

	// exec command `go fmt`
	err = utils.ExecCommand("./", "go", "mod", "tidy")
	if err != nil {
		return err
	}

	return nil
}
