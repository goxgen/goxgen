// Generated code will be used in runtime.
package main

import (
	"context"
	"fmt"
	"github.com/goxgen/goxgen/projects/basic"
	"github.com/goxgen/goxgen/projects/gorm"
	"github.com/goxgen/goxgen/xgen"
)

func main() {
	xg := xgen.NewXgen(
		xgen.WithPackageName("github.com/goxgen/goxgen/runtime"),
		xgen.WithProject(
			"simple_initial",
			basic.NewProject(),
		),
		xgen.WithProject(
			"gorm_initial",
			gorm.NewProject(),
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
