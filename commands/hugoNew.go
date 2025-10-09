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
	filename := slug.Make(title) + ".md"
	filepath := path.Join("content", "posts", filename)
	out, err := Execute("hugo", "new", "content", filepath)
	utils.OpenFileInEditor(path.Join(utils.HugoProject, filepath))
	return out, err
}

func Execute(command string, args ...string) ([]byte, error) {
	cmd := exec.Command(command, args...)
	cmd.Dir = utils.HugoProject // ← your Hugo project directory
	out, err := cmd.Output()
	if err != nil {
		return out, err
	}
	return out, nil
}
