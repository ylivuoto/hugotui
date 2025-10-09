// Package commands implements 'hugo new' commands
package commands

import (
	"fmt"
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

func Publish() ([]byte, []byte) {
	// FIX: build path, needs to be hugo project folder
	// TODO: port for scp
	build, buildError := Execute("hugo")
	upload, uploadError := Execute("scp", "-r", path.Join(utils.HugoProject, "public")+"/*", utils.HugoRemote)

	if buildError != nil {
		fmt.Println("Build error:", buildError)
	}
	if uploadError != nil {
		fmt.Println("Upload error:", uploadError)
	}

	return build, upload
}

func Preview() {
	// FIX: process won't stop on exit, kill
	go func() {
		out, error := Execute("hugo", "server")
		if error != nil {
			fmt.Println("Error:", error)
		}
		fmt.Println("Output:", string(out))
	}()
	Execute("xdg-open", "http://localhost:1313")
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
