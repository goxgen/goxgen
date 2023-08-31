package main

import (
	"context"
	"fmt"
	"github.com/goxgen/goxgen/plugins/cli"
	"github.com/goxgen/goxgen/projects/gorm"
	"github.com/goxgen/goxgen/projects/simple"
	"github.com/goxgen/goxgen/xgen"
)

func main() {
	xg := xgen.NewXgen(
		xgen.WithPlugin(cli.NewPlugin()),
		xgen.WithPackageName("github.com/goxgen/goxgen/internal/integration"),
		xgen.WithProject(
			"myproject",
			simple.NewPlugin(),
		),
		xgen.WithProject(
			"gormproj",
			gorm.NewPlugin(),
		),
		//xgen.WithProject(
		//	"entproj",
		//	projects.NewEntProject(),
		//),
	)

	err := xg.Generate(context.Background())
	if err != nil {
		fmt.Println(err)
	}
}
