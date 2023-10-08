package common

import "github.com/urfave/cli/v2"

const (
	ProjectFlag = "project"
	TestFlag    = "test"

	LogLevelFlag                 = "log-level"
	DevModeFlag                  = "dev-mode"
	HttpsFlag                    = "https"
	DatabaseDriverFlag           = "db-driver"
	DatabaseDsnFlag              = "db-dsn"
	HostFlag                     = "host"
	PortFlag                     = "port"
	AppPathFlag                  = "app-path"
	AppServerEnabledFlag         = "app-server-enabled"
	GraphqlUrlFlag               = "graphql-url"
	GraphqlUriPathFlag           = "graphql-uri-path"
	GraphqlPlaygroundUriPathFlag = "graphql-playground-uri-path"
	GraphqlPlaygroundEnabledFlag = "graphql-playground-enabled"
)

var (
	DynamicFlags = []string{
		LogLevelFlag,
		DevModeFlag,
		HttpsFlag,
		DatabaseDriverFlag,
		DatabaseDsnFlag,
		HostFlag,
		PortFlag,
		AppPathFlag,
		AppServerEnabledFlag,
		GraphqlUrlFlag,
		GraphqlUriPathFlag,
		GraphqlPlaygroundUriPathFlag,
		GraphqlPlaygroundEnabledFlag,
	}
)

type ProjectSettings struct {
	LogLevel                 string
	DevMode                  bool
	Https                    bool
	DatabaseDriver           string
	DatabaseDsn              string
	Host                     string
	Port                     string
	AppPath                  string
	AppServerEnabled         bool
	GraphqlUrl               string
	GraphqlUriPath           string
	GraphqlPlaygroundUriPath string
	GraphqlPlaygroundEnabled bool
}

func GetProjectSettings(c cli.Context) *ProjectSettings {
	return &ProjectSettings{
		LogLevel:                 c.String(LogLevelFlag),
		DevMode:                  c.Bool(DevModeFlag),
		Https:                    c.Bool(HttpsFlag),
		DatabaseDriver:           c.String(DatabaseDriverFlag),
		DatabaseDsn:              c.String(DatabaseDsnFlag),
		Host:                     c.String(HostFlag),
		Port:                     c.String(PortFlag),
		AppPath:                  c.String(AppPathFlag),
		AppServerEnabled:         c.Bool(AppServerEnabledFlag),
		GraphqlUrl:               c.String(GraphqlUrlFlag),
		GraphqlUriPath:           c.String(GraphqlUriPathFlag),
		GraphqlPlaygroundUriPath: c.String(GraphqlPlaygroundUriPathFlag),
		GraphqlPlaygroundEnabled: c.Bool(GraphqlPlaygroundEnabledFlag),
	}
}
