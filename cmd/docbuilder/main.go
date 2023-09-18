package main

import (
	"bytes"
	"fmt"
	"github.com/Masterminds/sprig/v3"
	"github.com/urfave/cli/v2"
	"go/ast"
	"go/parser"
	"go/token"
	"gopkg.in/yaml.v3"
	"io"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"sort"
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
						"funcSrc":       funcSrc,
						"fsTree":        fsTree,
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

func fsTree(root string, args ...any) string {
	treeStr := ""
	prefix := ""
	level := 2
	if len(args) > 0 {
		prefix = args[0].(string)
	}

	if len(args) > 1 {
		level = args[1].(int)
	}

	// Read the directory
	files, err := os.ReadDir(root)
	if err != nil {
		fmt.Println("Error reading directory:", err)
		return ""
	}

	// Sort the files, dirs first
	sort.Slice(files, func(i, j int) bool {
		if files[i].IsDir() && !files[j].IsDir() {
			return true
		}
		if !files[i].IsDir() && files[j].IsDir() {
			return false
		}
		return files[i].Name() < files[j].Name()
	})

	// Loop through each file/directory
	for i, file := range files {
		// Determine the new prefix for the next level
		newPrefix := prefix
		if i == len(files)-1 {
			treeStr += prefix + "└── " + file.Name() + "\n"
			newPrefix += "    "
		} else {
			treeStr += prefix + "├── " + file.Name() + "\n"
			newPrefix += "│   "
		}

		// If it's a directory, recurse
		if file.IsDir() && level > 0 {
			treeStr += fsTree(root+"/"+file.Name(), newPrefix, level-1)
		}
	}

	return treeStr
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

func funcSrc(fnNameWithPackage string) string {
	// Use regular expression to parse the function signature
	re := regexp.MustCompile(`^(.*?)(?:\.([^.]+))?\.(.+)$`)
	matches := re.FindStringSubmatch(fnNameWithPackage)
	if len(matches) != 4 {
		fmt.Println("Invalid function name with package")
		return ""
	}

	dirPath := matches[1]
	packageName := filepath.Base(dirPath)
	structName := matches[2]
	methodName := matches[3]

	fset := token.NewFileSet()

	var result string

	// Read all files in the specific directory
	files, err := os.ReadDir(dirPath)
	if err != nil {
		fmt.Println("Failed to read directory:", err)
		return ""
	}

	for _, file := range files {
		if file.IsDir() || filepath.Ext(file.Name()) != ".go" {
			continue
		}

		path := filepath.Join(dirPath, file.Name())
		src, err := os.ReadFile(path)
		if err != nil {
			continue
		}

		node, err := parser.ParseFile(fset, path, src, parser.ParseComments)
		if err != nil {
			fmt.Println("Failed to parse package:", err)
			continue
		}

		if node.Name.Name != packageName {
			continue
		}

		for _, decl := range node.Decls {
			if fn, isFn := decl.(*ast.FuncDecl); isFn {
				isMethodOfStruct := fn.Recv != nil && len(fn.Recv.List) > 0 &&
					fn.Recv.List[0].Type != nil && strings.Contains(fmt.Sprintf("%s", fn.Recv.List[0].Type), structName)

				isPackageLevelFunction := structName == "" && fn.Name.Name == methodName

				if isMethodOfStruct && fn.Name.Name == methodName || isPackageLevelFunction {
					// Extract the function source code using position information
					start := fset.Position(fn.Pos()).Offset
					end := fset.Position(fn.End()).Offset
					result = string(src[start:end])
					return result
				}
			}
		}
	}

	return result
}
