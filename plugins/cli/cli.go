package cli

import (
	"context"
	"embed"
	"github.com/goxgen/goxgen/plugins"
	"github.com/goxgen/goxgen/projects"
	"github.com/goxgen/goxgen/tmpl"
)

// CLI contains configuration of xgen-cli
type CLI struct {
}

func NewPlugin() *CLI {
	return &CLI{}
}

func (xc *CLI) Name() string {
	return "xgen-cli"
}

// templateData contains data for xgen-cli templates
type templateData struct {
	ParentPackageName string
	Projects          map[string]projects.Project
}

//go:embed templates/*
var templatesFS embed.FS

// Generate generates code from templates
func (xc *CLI) Generate(ctx context.Context) error {
	genCtx, err := plugins.GetContext(ctx)
	if err != nil {
		return err
	}

	data, err := xc.prepareCLITemplateData(genCtx)
	if err != nil {
		return err
	}

	tb := tmpl.TemplateBundle{
		TemplateDir: "templates",
		FS:          templatesFS,
		OutputFile:  "./" + genCtx.GeneratedFilePrefix + "cli.go",
		Regenerate:  true,
	}

	return tb.Generate(genCtx.OutputDir, data)
}

func (xc *CLI) prepareCLITemplateData(genCtx *plugins.Context) (*templateData, error) {
	return &templateData{
		Projects:          genCtx.Projects,
		ParentPackageName: genCtx.ParentPackageName,
	}, nil
}
