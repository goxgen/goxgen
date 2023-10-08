package common

import (
	"github.com/goxgen/goxgen/plugins/cli/consts"
	"github.com/goxgen/goxgen/plugins/cli/server"
	"github.com/goxgen/goxgen/plugins/cli/settings"
	"github.com/urfave/cli/v2"
	"io/fs"
)

type Project struct {
	Name    string
	Server  server.Constructor
	TestsFS fs.FS
	TestDir string
}

type ProjectList []*Project

type Settings struct {
	Projects *ProjectList
}

func NewProjectList(projects ...*Project) *ProjectList {
	pl := &ProjectList{}
	for _, p := range projects {
		pl.Add(p)
	}
	return pl
}

func (pl *ProjectList) Get(name string) *Project {
	for _, p := range *pl {
		if p.Name == name {
			return p
		}
	}
	return nil
}

func (pl *ProjectList) GetNames() []string {
	var names []string
	for _, p := range *pl {
		names = append(names, p.Name)
	}
	return names
}

func (pl *ProjectList) Add(p *Project) {
	*pl = append(*pl, p)
}

func (pl *ProjectList) ContainsAll(names ...string) bool {
	for _, name := range names {
		if pl.Get(name) == nil {
			return false
		}
	}
	return true
}

func prepareString(c *cli.Context, project string, flag string) string {
	projectFlag := FlagName(project, flag)
	if c.IsSet(projectFlag) {
		return c.String(projectFlag)
	}
	return c.String(flag)
}
func prepareBool(c *cli.Context, project string, flag string) bool {
	projectFlag := FlagName(project, flag)
	if c.IsSet(projectFlag) {
		return c.Bool(projectFlag)
	}
	return c.Bool(flag)
}
func prepareInt(c *cli.Context, project string, flag string) int {
	projectFlag := FlagName(project, flag)
	if c.IsSet(projectFlag) {
		return c.Int(projectFlag)
	}
	return c.Int(flag)
}

func (p *Project) GetEnvironmentSettings(c *cli.Context) *settings.EnvironmentSettings {
	return &settings.EnvironmentSettings{
		LogLevel:                 prepareString(c, p.Name, consts.LogLevelFlag),
		DevMode:                  prepareBool(c, p.Name, consts.DevModeFlag),
		Https:                    prepareBool(c, p.Name, consts.HttpsFlag),
		DatabaseDriver:           prepareString(c, p.Name, consts.DatabaseDriverFlag),
		DatabaseDsn:              prepareString(c, p.Name, consts.DatabaseDsnFlag),
		Host:                     prepareString(c, p.Name, consts.HostFlag),
		Port:                     prepareInt(c, p.Name, consts.PortFlag),
		AppPath:                  prepareString(c, p.Name, consts.AppPathFlag),
		AppServerEnabled:         prepareBool(c, p.Name, consts.AppServerEnabledFlag),
		GraphqlUrl:               prepareString(c, p.Name, consts.GraphqlUrlFlag),
		GraphqlUriPath:           prepareString(c, p.Name, consts.GraphqlUriPathFlag),
		GraphqlPlaygroundUriPath: prepareString(c, p.Name, consts.GraphqlPlaygroundUriPathFlag),
		GraphqlPlaygroundEnabled: prepareBool(c, p.Name, consts.GraphqlPlaygroundEnabledFlag),
	}
}
