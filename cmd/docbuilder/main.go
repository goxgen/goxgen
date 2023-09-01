package main

import (
	"bytes"
	"fmt"
	"github.com/Masterminds/sprig/v3"
	"github.com/urfave/cli/v2"
	"io"
	"os"
	"path/filepath"
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
					Funcs(template.FuncMap{"embed": embedFile}).
					Funcs(sprig.FuncMap())
				tmpl, err := tmpl.ParseFiles(templateFile)
				if err != nil {
					return err
				}

				f, err := os.Create(outputFile)
				if err != nil {
					return err
				}

				err = tmpl.Execute(f, nil)

				return nil
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Println(err.Error())
	}

}

func embedFile(filePath string) string {
	fileExtension := strings.TrimPrefix(filepath.Ext(filePath), ".")
	f, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	var buf bytes.Buffer
	if _, err := io.Copy(&buf, f); err != nil {
		panic(err)
	}

	return "```" + fileExtension + "\n" + buf.String() + "\n```"
}
