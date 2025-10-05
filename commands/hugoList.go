package commands

import (
	"encoding/csv"
	"os"
	"os/exec"
	"strings"
)

type Post struct {
	Path  string
	Slug  string
	Title string
	Date  string
}

func ListHugoPosts() ([]Post, error) {
	// Run Hugo command
	cmd := exec.Command("hugo", "list", "all")
	cmd.Dir = os.Getenv("HUGO_PATH") // ← your Hugo project directory
	out, err := cmd.Output()
	if err != nil {
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
	var posts []Post
	headers := records[0]
	for _, row := range records[1:] {
		p := Post{}
		for i, v := range row {
			switch headers[i] {
			case "path":
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
