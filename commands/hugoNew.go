// Package commands implements 'hugo new' commands
package commands

import (
	"os/exec"
	"path"

	"hugotui/utils"

	"github.com/gosimple/slug"
)

// CreateArticle implements 'hugo new content'
func CreateArticle(title string, tags []string) ([]byte, error) {
	// Run Hugo command
	// TODO: set path properly via input
	filename := slug.Make(title) + ".md"
	filepath := path.Join("content", "posts", filename)
	cmd := exec.Command("hugo", "new", "content", filepath)
	cmd.Dir = utils.HugoProject // ← your Hugo project directory
	out, err := cmd.Output()
	if err != nil {
		return out, err
	}

	// TODO: open recently created file
	utils.OpenFileInEditor(path.Join(cmd.Dir, filepath))

	return out, nil
}
