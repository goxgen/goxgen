package templates_engine

import (
	"embed"
	"fmt"
	"github.com/Masterminds/sprig/v3"
	"os"
	"path"
	"text/template"
)

// TemplateBundle contains configuration of a template bundle
type TemplateBundle struct {
	TemplateDir string           // template directory in the template bundle
	OutputFile  string           // output file
	Regenerate  bool             // regenerate output file if it already exists
	FS          embed.FS         //	template bundle file system
	FuncMap     template.FuncMap // template functions
}

// getTemplateNames returns template names in templates directory of the template bundle
func (tb *TemplateBundle) getTemplateNames() ([]string, error) {
	dirEntries, err := tb.FS.ReadDir(tb.TemplateDir)
	if err != nil {
		fmt.Println("Failed to read template directory:", err)
		return nil, err
	}

	var templateNames []string
	for _, entry := range dirEntries {
		if !entry.IsDir() {
			templateNames = append(templateNames, path.Join(tb.TemplateDir, entry.Name()))
		}
	}

	return templateNames, nil
}

func (tb *TemplateBundle) Generate(outputDir string, data any) error {
	fmt.Println("Generating:", *tb)
	outputFile := path.Join(outputDir, tb.OutputFile)
	var err error

	// check if output file exists and skip if not regenerate
	if !tb.Regenerate {
		if _, err := os.Stat(outputFile); !os.IsNotExist(err) {
			fmt.Println("Output file already exists, skipping:", outputFile)
			return nil
		}
	}

	templateNames, err := tb.getTemplateNames()
	if err != nil {
		fmt.Println("Failed to get template names:", err)
		return err
	}

	baseTemplate := template.
		New(tb.TemplateDir).
		Funcs(sprig.FuncMap()).
		Funcs(template.FuncMap{
			"indirect": Indirect,
		}).
		Funcs(tb.FuncMap)
	baseTemplate, err = baseTemplate.ParseFS(tb.FS, templateNames...)
	if err != nil {
		fmt.Println("Failed to parse templates:", err)
		return err
	}

	// create base directory if not exists
	err = os.MkdirAll(path.Dir(outputFile), 0755)
	if err != nil {
		fmt.Println("Failed to create base directory:", err)
		return err
	}

	file, err := os.Create(outputFile)
	if err != nil {
		fmt.Println("Failed to create output file:", err)
		return err
	}
	defer file.Close()

	err = baseTemplate.ExecuteTemplate(file, "BaseTemplate", data)
	if err != nil {
		fmt.Printf("Failed to execute template on file %s: %s\n", outputFile, err)
		return err
	}

	return nil
}
