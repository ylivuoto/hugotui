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
	utils.OpenFileInEditor(filepath)
	return out, err
}

func Publish() ([]byte, []byte) {
	// FIX: build path, needs to be hugo project folder
	// TODO: port for scp
	build, buildError := Execute("hugo")

	if buildError != nil {
		fmt.Println("Build error:", buildError)
	}
	go func() {
		upload, uploadError := Execute("bash", "-c", fmt.Sprintf("scp -r -P %s public/* %s", utils.HugoRemotePort, utils.HugoRemote))
		fmt.Println("Upload output:", string(upload))
		if uploadError != nil {
			fmt.Println(string(upload))
			fmt.Println("Upload error:", uploadError)
		}
	}()
	return build, nil
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
	out, err := cmd.Output()
	if err != nil {
		return out, err
	}
	return out, nil
}
