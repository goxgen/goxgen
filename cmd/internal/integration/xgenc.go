//go:build ignore
// +build ignore

package main

import (
	"context"
	"fmt"
	"github.com/goxgen/goxgen/plugins/cli"
	"github.com/goxgen/goxgen/projects/basic"
	"github.com/goxgen/goxgen/projects/gorm"
	"github.com/goxgen/goxgen/xgen"
)

func main() {
	xg := xgen.NewXgen(
		xgen.WithPackageName("github.com/goxgen/goxgen/cmd/internal/integration"),
		xgen.WithProject(
			"myproject",
			basic.NewProject(),
		),
		xgen.WithProject(
			"gorm_advanced",
			gorm.NewProject(
				gorm.WithBasicProjectOption(basic.WithTestDir("tests")),
			),
		),
		xgen.WithProject(
			"gorm_example",
			gorm.NewProject(
				gorm.WithBasicProjectOption(basic.WithTestDir("tests")),
			),
		),
		xgen.WithPlugin(cli.NewPlugin()),
	)

	err := xg.Generate(context.Background())
	if err != nil {
		fmt.Println(err)
	}
}
