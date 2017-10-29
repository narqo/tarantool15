package main

import (
	"bytes"
	"flag"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

var update = flag.Bool("update_golden", false, "update golden files")

func TestMain(m *testing.M) {
	flag.Parse()
	os.Exit(m.Run())
}

func TestGolden(t *testing.T) {
	cases := []struct {
		name  string
		in    string
		types []string
	}{
		{"token", "token.go", []string{"Token", "Service"}},
		{"private", "private.go", []string{"privateToken"}},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var g Generator
			inFiles := []string{filepath.Join("fixtures", tc.in)}
			g.parsePackageFiles(inFiles)
			for _, typeName := range tc.types {
				g.generate(typeName)
			}

			got := g.format()

			golden := filepath.Join("fixtures", tc.name+".golden")
			if *update {
				ioutil.WriteFile(golden, got, 0644)
			}
			want, err := ioutil.ReadFile(golden)
			if err != nil {
				t.Fatal(err)
			}
			if !bytes.Equal(got, want) {
				t.Errorf("%s: got\n====\n%s\n====\nwant\n====\n%s", tc.name, got, want)
			}
		})
	}
}
