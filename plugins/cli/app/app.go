package common

import (
	"fmt"
	"github.com/goxgen/goxgen/plugins/cli/consts"
	"github.com/goxgen/goxgen/plugins/cli/project"
	"github.com/goxgen/goxgen/plugins/cli/server"
	"github.com/goxgen/goxgen/plugins/cli/tests"
	"github.com/joho/godotenv"
	"github.com/urfave/cli/v2"
	"os"
	"sync"
)

func newApp(settings *project.Settings) *cli.App {
	app := cli.NewApp()
	app.Name = "GoXGen"
	app.Version = "0.1.0"
	app.Description = "This is GoXGen CLI"
	app.Authors = []*cli.Author{
		{
			Name:  "Aaron Yordanyan",
			Email: "aaron.yor@gmail.com",
		},
	}

	app.Commands = []*cli.Command{
		GenerateRunCommand(settings.Projects),
	}
	return app
}

func GenerateRunCommand(projects *project.ProjectList) *cli.Command {
	flags := append([]cli.Flag{
		&cli.StringSliceFlag{
			Name:     consts.ProjectFlag,
			Required: true,
		},
	})

	for _, project := range append([]string{""}, projects.GetNames()...) {
		flags = append(flags,
			&cli.BoolFlag{
				Name:    FlagName(project, consts.TestFlag),
				EnvVars: []string{envName(project, consts.TestFlag)},
			},
			&cli.StringFlag{
				Name:    FlagName(project, consts.LogLevelFlag),
				Value:   "debug",
				EnvVars: []string{envName(project, consts.LogLevelFlag)},
			},
			&cli.BoolFlag{
				Name:    FlagName(project, consts.DevModeFlag),
				Value:   true,
				EnvVars: []string{envName(project, consts.DevModeFlag)},
			},
			&cli.BoolFlag{
				Name:    FlagName(project, consts.HttpsFlag),
				Value:   false,
				EnvVars: []string{envName(project, consts.HttpsFlag)},
			},
			&cli.StringFlag{
				Name:    FlagName(project, consts.DatabaseDriverFlag),
				Value:   "sqlite",
				EnvVars: []string{envName(project, consts.DatabaseDriverFlag)},
			},
			&cli.StringFlag{
				Name:    FlagName(project, consts.DatabaseDsnFlag),
				Value:   "file:demo.db?mode=rwc&cache=shared&_fk=1",
				EnvVars: []string{envName(project, consts.DatabaseDsnFlag)},
			},
			&cli.StringFlag{
				Name:    FlagName(project, consts.AppPathFlag),
				Value:   "/",
				EnvVars: []string{envName(project, consts.AppPathFlag)},
			},
			&cli.BoolFlag{
				Name:    FlagName(project, consts.AppServerEnabledFlag),
				Value:   true,
				EnvVars: []string{envName(project, consts.AppServerEnabledFlag)},
			},
			&cli.StringFlag{
				Name:    FlagName(project, consts.GraphqlUrlFlag),
				EnvVars: []string{envName(project, consts.GraphqlUrlFlag)},
			},
			&cli.StringFlag{
				Name:    FlagName(project, consts.GraphqlUriPathFlag),
				Value:   "/query",
				EnvVars: []string{envName(project, consts.GraphqlUriPathFlag)},
			},
			&cli.StringFlag{
				Name:    FlagName(project, consts.GraphqlPlaygroundUriPathFlag),
				Value:   "/playground",
				EnvVars: []string{envName(project, consts.GraphqlPlaygroundUriPathFlag)},
			},
			&cli.BoolFlag{
				Name:    FlagName(project, consts.GraphqlPlaygroundEnabledFlag),
				Value:   false,
				EnvVars: []string{envName(project, consts.GraphqlPlaygroundEnabledFlag)},
			},
			&cli.StringFlag{
				Name:    FlagName(project, consts.HostFlag),
				Value:   "localhost",
				EnvVars: []string{envName(project, consts.HostFlag)},
			},
			&cli.IntFlag{
				Name:    FlagName(project, consts.PortFlag),
				Value:   80,
				EnvVars: []string{envName(project, consts.PortFlag)},
			},
		)
	}

	cmd := &cli.Command{
		Name:  "run",
		Flags: flags,
		Action: func(c *cli.Context) error {
			projectNames := c.StringSlice(consts.ProjectFlag)
			test := c.Bool(consts.TestFlag)

			if !projects.ContainsAll(projectNames...) {
				return fmt.Errorf("invalid project name")
			}

			wg := sync.WaitGroup{}
			for _, projectName := range projectNames {
				project := projects.Get(projectName)
				sts := project.GetEnvironmentSettings(c)

				if test {
					return tests.Start(sts, project)
				} else {
					wg.Add(1)
					srv, err := server.New(sts)
					if err != nil {
						return err
					}
					go func() {
						err := srv.ListenAndServe(project.Server)
						defer wg.Done()
						if err != nil {
							fmt.Println(err.Error())
						}
					}()
				}
			}

			if !test {
				wg.Wait()
			}
			return nil
		},
	}
	return cmd
}

func Run(settings *project.Settings) {
	app := newApp(settings)
	_ = godotenv.Load(".env.local", ".env")
	if err := app.Run(os.Args); err != nil {
		fmt.Println(err.Error())
	}
}
