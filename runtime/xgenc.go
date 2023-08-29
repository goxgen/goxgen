// Generated code will be used in runtime.
package main

import (
	"context"
	"fmt"
	"github.com/goxgen/goxgen/projects/gorm"
	"github.com/goxgen/goxgen/projects/simple"
	"github.com/goxgen/goxgen/xgen"
)

func main() {
	xg := xgen.NewXgen(
		xgen.WithPackageName("github.com/goxgen/goxgen/runtime"),
		xgen.WithProject(
			"simple_initial",
			simple.New(),
		),
		xgen.WithProject(
			"gorm_initial",
			gorm.New(),
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
