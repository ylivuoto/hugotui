package utils

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"strings"

	"github.com/gosimple/slug"
)

// ReadFileAsString reads the file at the given path and returns its contents as a string.
func ReadFileAsString(path string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func OpenFileInEditor(filePath string) error {
	editor := os.Getenv("EDITOR")
	terminal := os.Getenv("TERMINAL")
	if editor == "" {
		editor = "nvim"
	}
	cmd := exec.Command(terminal, "-e", editor, filePath)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// m.status = fmt.Sprintf("Opening %s...", filePath)
	go func() {
		if err := cmd.Run(); err != nil {
			fmt.Printf("Failed: %v", err)
		}
	}()
	// else {
	// 	m.status = fmt.Sprintf("Closed editor for %s", filePath)
	// }
	// TODO: return error properly
	return nil
}

func ModifyFileTitle(filepath string, title string) error {
	fullPath := filepath
	data, err := os.ReadFile(fullPath)
	if err != nil {
		return err
	}

	lines := strings.Split(string(data), "\n")
	for i, line := range lines {
		if strings.HasPrefix(line, "title = ") {
			lines[i] = fmt.Sprintf("title =  \"%s\"", title)
			break
		}
	}
	return os.WriteFile(fullPath, []byte(strings.Join(lines, "\n")), 0o644)
}

func ModifyFileTags(filepath string, tags []string) error {
	fullPath := filepath
	data, err := os.ReadFile(fullPath)
	if err != nil {
		return err
	}

	lines := strings.Split(string(data), "\n")
	for i, line := range lines {
		if strings.HasPrefix(line, "tags = ") {
			lines[i] = fmt.Sprintf("tags = [\"%s\"]", strings.Join(tags, "\", \""))
			break
		}
	}
	return os.WriteFile(fullPath, []byte(strings.Join(lines, "\n")), 0o644)
}

// ModifyFilePath renames the file based on the new title and moves it to the posts directory.
func ModifyFilePath(filepath string, title string) error {
	filename := slug.Make(title) + ".md"
	newPath := path.Join("content", "posts", filename)
	return os.Rename(filepath, newPath)
}
