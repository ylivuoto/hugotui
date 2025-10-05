// Package utils provides utility functions for parsing front matter from content strings.
package utils

import (
	"strings"

	"github.com/adrg/frontmatter"
)

// Struct that works for both TOML and YAML
var matter struct {
	Title string   `yaml:"title"`
	Tags  []string `yaml:"tags"`
	Date  string   `yaml:"date"`
}

type Matter struct {
	Title string   `yaml:"title"`
	Tags  []string `yaml:"tags"`
	Date  string   `yaml:"date"`
}

func ParseFrontMatter(content string) (Matter, string) {
	rest, err := frontmatter.Parse(strings.NewReader(content), &matter)
	if err != nil {
		println("Error parsing front matter:", err.Error())
	}
	// NOTE: If a front matter must be present in the input data, use
	//       frontmatter.MustParse instead.

	// fmt.Printf("%+v\n", matter)
	// fmt.Println(string(rest))

	// Output:
	// {Name:frontmatter Tags:[go yaml json toml]}
	// rest of the content
	//
	return Matter(matter), string(rest)
}
