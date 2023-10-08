package app

import (
	"fmt"
	"github.com/goxgen/goxgen/plugins/cli/project"
	"github.com/joho/godotenv"
	"github.com/urfave/cli/v2"
	"os"
)

type Settings struct {
	Projects *project.List
}

// newApp creates a new cli app with urfave/cli based on provided settings
func newApp(settings *Settings) *cli.App {
	app := cli.NewApp()
	app.Name = "goxgen"
	app.Version = "0.1.0"
	app.Description = "This is goxgen cli"
	app.Authors = []*cli.Author{
		{
			Name: "goxgen",
		},
		{
			Name:  "Aaron Yordanyan",
			Email: "aaron.yor@gmail.com",
		},
	}

	app.Commands = []*cli.Command{GenerateRunCommand(settings.Projects)}
	return app
}

// Run runs the cli app
func Run(settings *Settings) {
	app := newApp(settings)

	_ = godotenv.Load(".env", ".env.default")
	_ = godotenv.Overload(".env.local")

	if err := app.Run(os.Args); err != nil {
		fmt.Println(err.Error())
	}
}
