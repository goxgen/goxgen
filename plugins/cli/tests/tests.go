package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/goxgen/goxgen/plugins/cli/server"
	"github.com/urfave/cli/v2"
	"golang.org/x/exp/slices"
	"gopkg.in/yaml.v3"
	"io"
	"io/fs"
	"net/http"
	"path/filepath"
	"strconv"
)

type TestBundle struct {
	Name  string      `yaml:"name"`
	Tests []*TestCase `yaml:"tests"`
}

type TestCase struct {
	Name           string `yaml:"name"`
	Query          string `yaml:"query"`
	ExpectedResult string `yaml:"expectedResult"`
	Priority       int    `yaml:"priority"`
}

// pureJson removes all whitespaces from a json string
func pureJson(jsonStr string) string {
	dst := &bytes.Buffer{}
	if err := json.Compact(dst, []byte(jsonStr)); err != nil {
		panic(err)
	}
	return dst.String()
}

// Run runs a single test case
func (b *TestCase) Run(testClient *http.Client, gqlEndpoint string) error {
	fmt.Println("--> Testcase " + b.Name)
	body := map[string]any{}
	body["query"] = b.Query
	bodyJson, err := json.Marshal(body)
	if err != nil {
		return err
	}

	resp, err := testClient.Post(
		gqlEndpoint,
		"application/json",
		bytes.NewReader(bodyJson),
	)

	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			fmt.Println("Bad content", err)
		}
		bodyString := string(bodyBytes)
		gotResult := pureJson(bodyString)
		expectedResult := pureJson(b.ExpectedResult)
		if gotResult == expectedResult {
			return nil
		} else {
			return fmt.Errorf("failed test: expected %v, got %v", expectedResult, gotResult)
		}
	}

	return fmt.Errorf("bad response %v", resp)
}

// Run runs all tests in a bundle
func (b *TestBundle) Run(testClient *http.Client, gqlEndpoint string) error {
	testsCount := len(b.Tests)
	fmt.Println("-> Running " + b.Name + " (" + strconv.Itoa(testsCount) + ")")

	for _, tc := range b.Tests {
		err := tc.Run(testClient, gqlEndpoint)
		if err != nil {
			return err
		}
	}
	return nil
}

// find finds all files with given extensions in a directory
func find(testsFs fs.FS, root string, extensions ...string) []string {
	var a []string
	err := fs.WalkDir(testsFs, root, func(s string, d fs.DirEntry, e error) error {
		if e != nil {
			return e
		}
		if slices.Contains(extensions, filepath.Ext(d.Name())) {
			a = append(a, s)
		}
		return nil
	})
	if err != nil {
		panic(err)
	}
	return a
}

// getTestBundle reads a test bundle from a file
func getTestBundle(testsFS fs.FS, file string) (*TestBundle, error) {
	f, err := testsFS.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	var bundle TestBundle
	err = yaml.NewDecoder(f).Decode(&bundle)
	if err != nil {
		return nil, err
	}
	return &bundle, nil
}

// getTestBundles reads all test bundles from a directory
func getTestBundles(testsFS fs.FS, dir string) []*TestBundle {
	files := find(testsFS, dir, ".yaml", ".yml")
	var bundles []*TestBundle
	for _, file := range files {
		bundle, err := getTestBundle(testsFS, file)
		if err != nil {
			panic(err)
		}
		bundles = append(bundles, bundle)
	}
	return bundles
}

// Start runs all tests
func Start(ctx *cli.Context, serverConstructor server.Constructor, testsFS fs.FS, testsDirectory string) error {

	srv, err := server.New(ctx)
	if err != nil {
		return err
	}

	srvData := srv.GetDataFromCliContext()

	testSrv, cancel := srv.TestServer(ctx, serverConstructor)
	defer cancel()

	testClient := testSrv.Client()

	gqlEndpoint := testSrv.URL + srvData.GraphqlURIPath

	tbs := getTestBundles(testsFS, testsDirectory)
	for _, tb := range tbs {

		err := tb.Run(testClient, gqlEndpoint)
		if err != nil {
			return err
		}
	}

	return nil
}
