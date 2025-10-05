// Package commands implements 'hugo new' commands
package commands

import (
	"os"
	"os/exec"
	"path"
)

// CreateArticle implements 'hugo new content'
func CreateArticle() error {
	// Run Hugo command
	// TODO: set path properly via input
	path := path.Join("content", "posts", "temp.md")
	cmd := exec.Command("hugo", "new", "content", path)
	cmd.Dir = os.Getenv("HUGO_PATH") // ← your Hugo project directory
	out, err := cmd.Output()
	if err != nil {
		return err
	}
	// TODO: open recently created file
	println(string(out))

	return nil
}
