package goxgen

import (
	"context"
)

// cliTemplateBundle contains template bundle for xgen-cli
var cliTemplateBundle TemplateBundle = TemplateBundle{
	TemplateDir: "templates/xgen-cli-templates",
	FS:          templatesFS,
	OutputFile:  "./generated_main.go",
	Regenerate:  true,
}

// CLI contains configuration of xgen-cli
type CLI struct {
	OutputDir         *string
	Projects          *[]*SimpleProject
	ParentPackageName *string
}

type CLITemplateData struct {
	ParentPackageName string
	Projects          []Project
}

func (xc *CLI) prepareCLITemplateData(ctx context.Context) (*CLITemplateData, error) {
	genCtx, err := GetGeneratorContext(ctx)
	if err != nil {
		return nil, err
	}

	return &CLITemplateData{
		Projects:          genCtx.Projects,
		ParentPackageName: genCtx.ParentPackageName,
	}, nil
}

// Generate generates code from templates
func (xc *CLI) Generate(ctx context.Context) error {

	data, err := xc.prepareCLITemplateData(ctx)
	if err != nil {
		return err
	}

	return cliTemplateBundle.generate(PString(xc.OutputDir), data)
}
