// Package commands: implements 'hugo list' commands
package commands

import (
	"encoding/csv"
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"

	"hugotui/utils"
)

type Post struct {
	Path    string
	Slug    string
	Title   string
	Date    string
	Content string
	Tags    []string
}

// ListHugoPosts implement 'hugo list all' and parse CSV output
func ListHugoPosts() ([]Post, error) {
	// Run Hugo command
	cmd := exec.Command("hugo", "list", "all")
	out, err := cmd.Output()
	if err != nil {
		fmt.Println("Error executing hugo list in: ", cmd.Dir)
		return nil, err
	}

	// Read CSV
	reader := csv.NewReader(strings.NewReader(string(out)))
	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	if len(records) < 2 {
		return nil, nil // no data
	}

	// Parse each record
	// TODO: refactor name everywhere, article?
	numPosts := len(records) - 1
	posts := make([]Post, numPosts)

	for i, row := range records[1:] {
		p := Post{}
		content, _ := utils.ReadFileAsString(filepath.Join(cmd.Dir, row[0]))
		matter, body := utils.ParseFrontMatter(content)
		fmt.Println("Front matter: ", matter)
		p.Tags = append([]string(nil), matter.Tags...)
		p.Content = body
		p.Path = row[0]
		p.Slug = row[1]
		p.Title = strings.Trim(row[2], `"`)
		p.Date = row[3]
		posts[i] = p
	}

	return posts, nil
}
