package settings

type EnvironmentSettings struct {
	LogLevel                 string
	DevMode                  bool
	Https                    bool
	DatabaseDriver           string
	DatabaseDsn              string
	Host                     string
	Port                     int
	AppPath                  string
	AppServerEnabled         bool
	GraphqlUrl               string
	GraphqlUriPath           string
	GraphqlPlaygroundUriPath string
	GraphqlPlaygroundEnabled bool
}
