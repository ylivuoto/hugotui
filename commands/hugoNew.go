// Package commands implements 'hugo new' commands
package commands

import (
	"fmt"
	"os"
	"os/exec"
	"path"

	"hugotui/utils"

	"github.com/gosimple/slug"
)

var hugoProcess *os.Process

// CreateArticle implements 'hugo new content'
func CreateArticle(title string, tags []string) ([]byte, error) {
	// Run Hugo command
	filename := slug.Make(title) + ".md"
	filepath := path.Join("content", "posts", filename)
	out, err := Execute("hugo", "new", "content", filepath)

	// Run editor in goroutine for parallelism
	utils.OpenFileInEditor(filepath)
	return out, err
}

func Publish() ([]byte, []byte) {
	// TODO: refactor publish, this function is out of use
	build, buildError := Execute("hugo")

	if buildError != nil {
		fmt.Println("Build error:", buildError)
	}
	go func() {
		upload, uploadError := Execute("bash", "-c", fmt.Sprintf("scp -r public/* %s:%s", utils.HugoRemote, utils.HugoRemoteDir))
		fmt.Println("Upload output:", string(upload))
		if uploadError != nil {
			fmt.Println(string(upload))
			fmt.Println("Upload error:", uploadError)
		}
	}()
	return build, nil
}

func Preview() {
	cmd := exec.Command("hugo", "server")
	err := cmd.Start()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	hugoProcess = cmd.Process
	go func() {
		Execute("xdg-open", "http://localhost:1313")
	}()
}

// StopPreview calls stops existing hugo server process
func StopPreview() {
	if hugoProcess != nil {
		hugoProcess.Kill() // or hugoProcess.Signal(os.Interrupt)
		hugoProcess = nil
	}
}

func Execute(command string, args ...string) ([]byte, error) {
	cmd := exec.Command(command, args...)
	out, err := cmd.Output()
	if err != nil {
		return out, err
	}
	return out, nil
}
