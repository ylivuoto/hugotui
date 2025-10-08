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
	cmd.Dir = utils.HugoProject // ← your Hugo project directory
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
	var posts []Post
	headers := records[0]
	for _, row := range records[1:] {
		p := Post{}
		for i, v := range row {
			switch headers[i] {
			case "path":
				// TODO: handle error
				content, _ := utils.ReadFileAsString(filepath.Join(cmd.Dir, v))
				matter, body := utils.ParseFrontMatter(content)
				p.Tags = matter.Tags
				p.Content = body
				p.Path = v
			case "slug":
				p.Slug = v
			case "title":
				p.Title = strings.Trim(v, `"`)
			case "date":
				p.Date = v
			}
		}
		posts = append(posts, p)
	}
	return posts, nil
}
