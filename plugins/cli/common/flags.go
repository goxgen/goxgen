package common

import (
	"github.com/iancoleman/strcase"
	"strings"
)

func FlagName(project string, name string) string {
	str := strings.Join([]string{project, name}, "_")
	str = strings.Trim(str, "_")
	return str
}

func EnvName(project string, name string) string {
	flagName := FlagName(project, name)
	return strcase.ToScreamingSnake(flagName)
}
