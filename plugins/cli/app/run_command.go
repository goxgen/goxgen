package app

import (
	"fmt"
	"github.com/goxgen/goxgen/plugins/cli/common"
	"github.com/goxgen/goxgen/plugins/cli/consts"
	"github.com/goxgen/goxgen/plugins/cli/project"
	"github.com/goxgen/goxgen/plugins/cli/server"
	"github.com/goxgen/goxgen/plugins/cli/tests"
	"github.com/urfave/cli/v2"
	"sync"
)

// GenerateRunCommand generates run command based on projects
func GenerateRunCommand(projects *project.List) *cli.Command {
	flags := append([]cli.Flag{
		&cli.StringSliceFlag{
			Name:        consts.ProjectFlag,
			Usage:       "projects to run",
			Value:       cli.NewStringSlice(projects.GetNames()...),
			DefaultText: fmt.Sprintf("Available projects: %v", projects.GetNames()),
		},
	})

	for _, pj := range append([]string{""}, projects.GetNames()...) {
		flags = append(flags,
			&cli.BoolFlag{
				Name:    common.FlagName(pj, consts.TestFlag),
				EnvVars: []string{common.EnvName(pj, consts.TestFlag)},
			},
			&cli.StringFlag{
				Name:    common.FlagName(pj, consts.LogLevelFlag),
				Value:   "debug",
				EnvVars: []string{common.EnvName(pj, consts.LogLevelFlag)},
			},
			&cli.BoolFlag{
				Name:    common.FlagName(pj, consts.DevModeFlag),
				Value:   true,
				EnvVars: []string{common.EnvName(pj, consts.DevModeFlag)},
			},
			&cli.BoolFlag{
				Name:    common.FlagName(pj, consts.HttpsFlag),
				Value:   false,
				EnvVars: []string{common.EnvName(pj, consts.HttpsFlag)},
			},
			&cli.StringFlag{
				Name:    common.FlagName(pj, consts.DatabaseDriverFlag),
				Value:   "sqlite",
				EnvVars: []string{common.EnvName(pj, consts.DatabaseDriverFlag)},
			},
			&cli.StringFlag{
				Name:    common.FlagName(pj, consts.DatabaseDsnFlag),
				Value:   "file:demo.db?mode=rwc&cache=shared&_fk=1",
				EnvVars: []string{common.EnvName(pj, consts.DatabaseDsnFlag)},
			},
			&cli.StringFlag{
				Name:    common.FlagName(pj, consts.AppPathFlag),
				Value:   "/",
				EnvVars: []string{common.EnvName(pj, consts.AppPathFlag)},
			},
			&cli.BoolFlag{
				Name:    common.FlagName(pj, consts.AppServerEnabledFlag),
				Value:   true,
				EnvVars: []string{common.EnvName(pj, consts.AppServerEnabledFlag)},
			},
			&cli.StringFlag{
				Name:    common.FlagName(pj, consts.GraphqlUrlFlag),
				EnvVars: []string{common.EnvName(pj, consts.GraphqlUrlFlag)},
			},
			&cli.StringFlag{
				Name:    common.FlagName(pj, consts.GraphqlUriPathFlag),
				Value:   "/query",
				EnvVars: []string{common.EnvName(pj, consts.GraphqlUriPathFlag)},
			},
			&cli.StringFlag{
				Name:    common.FlagName(pj, consts.GraphqlPlaygroundUriPathFlag),
				Value:   "/playground",
				EnvVars: []string{common.EnvName(pj, consts.GraphqlPlaygroundUriPathFlag)},
			},
			&cli.BoolFlag{
				Name:    common.FlagName(pj, consts.GraphqlPlaygroundEnabledFlag),
				Value:   false,
				EnvVars: []string{common.EnvName(pj, consts.GraphqlPlaygroundEnabledFlag)},
			},
			&cli.StringFlag{
				Name:    common.FlagName(pj, consts.HostFlag),
				Value:   "localhost",
				EnvVars: []string{common.EnvName(pj, consts.HostFlag)},
			},
			&cli.IntFlag{
				Name:    common.FlagName(pj, consts.PortFlag),
				Value:   80,
				EnvVars: []string{common.EnvName(pj, consts.PortFlag)},
			},
		)
	}

	return &cli.Command{
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
				pj := projects.Get(projectName)
				pj.Init(c)

				if test {
					return tests.Start(pj)
				} else {
					wg.Add(1)
					srv, err := server.New(pj)
					if err != nil {
						return err
					}
					go func() {
						err := srv.ListenAndServe(pj.Server)
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
}
