package goxgen

import (
	"context"
	"github.com/goxgen/goxgen/templates_engine"
)

// CLI contains configuration of xgen-cli
type CLI struct {
	OutputDir         *string
	Projects          *[]*SimpleProject
	ParentPackageName *string
}

// CLITemplateData contains data for xgen-cli templates
type CLITemplateData struct {
	ParentPackageName string
	Projects          []Project
}

// cliTemplateBundle contains template bundle for xgen-cli
var cliTemplateBundle = templates_engine.TemplateBundle{
	TemplateDir: "templates/xgen-cli-templates",
	FS:          templatesFS,
	OutputFile:  "./generated_main.go",
	Regenerate:  true,
}

// Generate generates code from templates
func (xc *CLI) Generate(ctx context.Context) error {

	data, err := xc.prepareCLITemplateData(ctx)
	if err != nil {
		return err
	}

	return cliTemplateBundle.Generate(PString(xc.OutputDir), data)
}

func (xc *CLI) prepareCLITemplateData(ctx context.Context) (*CLITemplateData, error) {
	genCtx, err := GetXgenContext(ctx)
	if err != nil {
		return nil, err
	}

	return &CLITemplateData{
		Projects:          genCtx.Projects,
		ParentPackageName: genCtx.ParentPackageName,
	}, nil
}
