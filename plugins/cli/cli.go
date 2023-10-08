// Package cli contains the cli plugin that generates application entrypoint
// and cli commands for the projects: e.g. graphql server with playground.
// Also, it gets environment variables and cli arguments and injects them into the application.
// You can run multiple projects at the same time in one command or can run them separately.
package cli

import (
	"context"
	"embed"
	"github.com/goxgen/goxgen/consts"
	"github.com/goxgen/goxgen/plugins"
	"github.com/goxgen/goxgen/projects"
	"github.com/goxgen/goxgen/tmpl"
	"path"
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
	ParentPackageName          string
	GeneratedGqlgenPackageName string
	Projects                   map[string]projects.Project
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

	tbl := tmpl.TemplateBundleList{
		&tmpl.TemplateBundle{
			TemplateDir: "templates/main",
			FS:          templatesFS,
			OutputFile:  path.Join(genCtx.GeneratedFilePrefix + "cli.go"),
			Regenerate:  true,
		},
		&tmpl.TemplateBundle{
			TemplateDir: "templates/env",
			FS:          templatesFS,
			OutputFile:  path.Join(".env"),
			Regenerate:  false,
		},
		&tmpl.TemplateBundle{
			TemplateDir: "templates/env_default",
			FS:          templatesFS,
			OutputFile:  path.Join(".env.default"),
			Regenerate:  true,
		},
	}

	return tbl.Generate(genCtx.OutputDir, data)
}

func (xc *CLI) prepareCLITemplateData(genCtx *plugins.Context) (*templateData, error) {
	return &templateData{
		Projects:                   genCtx.Projects,
		ParentPackageName:          genCtx.ParentPackageName,
		GeneratedGqlgenPackageName: consts.GeneratedGqlgenPackageName,
	}, nil
}
