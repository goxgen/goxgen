package main

import (
	"bytes"
	"fmt"
	"github.com/Masterminds/sprig/v3"
	"github.com/urfave/cli/v2"
	"gopkg.in/yaml.v3"
	"io"
	"os"
	"reflect"
	"strings"
	"text/template"
)

func main() {

	app := cli.NewApp()
	app.Name = "GoXGen Document Builder"
	app.Version = "0.1.0"
	app.Description = "This is GoXGen CLI"
	app.Authors = []*cli.Author{
		{
			Name:  "Aaron Yordanyan",
			Email: "aaron.yor@gmail.com",
		},
	}

	app.Commands = []*cli.Command{
		{
			Name:    "build",
			Aliases: []string{"b"},
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:     "Template",
					Aliases:  []string{"template", "t"},
					Usage:    "Template to use for generation",
					Required: true,
				},
				&cli.StringFlag{
					Name:     "Output",
					Aliases:  []string{"output", "o"},
					Usage:    "Output file",
					Required: true,
				},
			},
			Action: func(c *cli.Context) error {
				templateFile := c.String("Template")
				outputFile := c.String("Output")

				// Implement the logic for building the documentation
				// Take the template file and rend with golang/text/template
				// Save the output to the output file
				tmpl := template.
					New(templateFile).
					Funcs(sprig.FuncMap()).
					Funcs(template.FuncMap{
						"embedFile":     embedFile,
						"codeBlock":     codeBlock,
						"yamlFileParse": yamlFileParse,
						"dig":           dig,
						"lines":         lines,
					})

				tmpl = template.Must(tmpl.ParseFiles(templateFile))

				f, err := os.Create(outputFile)
				if err != nil {
					return err
				}

				return tmpl.Execute(f, nil)
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Println(err.Error())
	}

}

func getFileContent(filePath string) (string, error) {
	f, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	var buf bytes.Buffer
	if _, err := io.Copy(&buf, f); err != nil {
		panic(err)
	}
	return buf.String(), nil
}

func embedFile(filePath string) (string, error) {
	content, err := getFileContent(filePath)
	if err != nil {
		return "", err
	}
	return content, nil
}

func codeBlock(language string, content string) string {
	return "```" + language + "\n" + content + "\n```"
}

func lines(args ...any) string {
	content, ok := args[len(args)-1].(string)
	if !ok {
		panic("lines needs a string as the last argument")
	}

	args = args[:len(args)-1]
	lines := strings.Split(content, "\n")
	totalLines := len(lines)
	lineStart := 0
	lineEnd := totalLines
	if len(args) > 0 {
		lineStart, ok = args[0].(int)
		if !ok {
			panic("lines needs an integer as the first argument")
		}
	}

	if len(args) > 1 {
		lineEnd, ok = args[1].(int)
		if !ok {
			panic("lines needs an integer as the second argument")
		}
	}

	return strings.Join(lines[lineStart:lineEnd], "\n")
}

func yamlFileParse(filePath string) (any, error) {
	content, err := getFileContent(filePath)
	if err != nil {
		return nil, err
	}
	var data any
	err = yaml.Unmarshal([]byte(content), &data)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func dig(ps ...any) (any, error) {
	if len(ps) < 3 {
		panic("dig needs at least three arguments")
	}
	dict := toMap(ps[len(ps)-1])
	def := ps[len(ps)-2]
	ks := make([]any, len(ps)-2)
	for i := 0; i < len(ks); i++ {
		ks[i] = ps[i]
	}

	return digFromDict(dict, def, ks)
}

func digFromDict(dict map[any]any, d any, ks []any) (any, error) {
	k, ns := ks[0], ks[1:]
	step, has := dict[k]
	if !has {
		return d, nil
	}
	if len(ns) == 0 {
		return step, nil
	}
	return digFromDict(toMap(step), d, ns)
}

func toMap(i any) map[any]any {
	reflectType := reflect.TypeOf(i)

	if reflectType.Kind() == reflect.Map {
		m := make(map[any]any)
		for _, key := range reflect.ValueOf(i).MapKeys() {
			m[key.Interface()] = reflect.ValueOf(i).MapIndex(key).Interface()
		}
		return m
	}

	if reflectType.Kind() == reflect.Ptr {
		return toMap(reflect.ValueOf(i).Elem().Interface())
	}

	if reflectType.Kind() == reflect.Slice {
		slice := reflect.ValueOf(i)
		m := make(map[any]any)
		for i := 0; i < slice.Len(); i++ {
			m[i] = slice.Index(i).Interface()
		}
		return m
	}

	panic("toMap error: not a map or slice")
}
