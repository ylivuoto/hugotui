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
	// TODO: refactor to reduce code duplication
	lines := getPostLines(filepath)
	modifyLines(&lines, title, "title = ", "title =  \"%s\"", "")
	return os.WriteFile(filepath, []byte(strings.Join(lines, "\n")), 0o644)
}

func ModifyFileTags(filepath string, tags []string) error {
	lines := getPostLines(filepath)
	modifyLines(&lines, tags, "tags = ", "tags = [\"%s\"]", "\", \"")
	return os.WriteFile(filepath, []byte(strings.Join(lines, "\n")), 0o644)
}

func ModifyExpiryDate(filepath string, date string) error {
	lines := getPostLines(filepath)
	for _, line := range lines {
		if strings.HasPrefix(strings.TrimSpace(line), "expiryDate = ") {
			return fmt.Errorf("expiryDate already set")
		}
	}

	expiryDate := fmt.Sprintf("expiryDate = \"%s\"", date)
	modifyLines(&lines, expiryDate, "expiryDate = ", "expiryDate = \"%s\"", "")
	return os.WriteFile(filepath, []byte(strings.Join(lines, "\n")), 0o644)
}

func modifyLines(lines *[]string, value any, prefix, format, separator string) {
	var content string
	switch val := value.(type) {
	case []string:
		content = strings.Join(val, separator)
	case string:
		content = val
	}

	exists := false
	for i, line := range *lines {
		if strings.HasPrefix(line, prefix) {
			(*lines)[i] = fmt.Sprintf(format, content)
			exists = true
			break
		}
	}

	if !exists {
		*lines = append((*lines)[:1], append([]string{content}, (*lines)[1:]...)...)
	}
}

func getPostLines(filepath string) []string {
	data, err := os.ReadFile(filepath)
	if err != nil {
		return []string{}
	}

	lines := strings.Split(string(data), "\n")

	return lines
}

// ModifyFilePath renames the file based on the new title and moves it to the posts directory.
func ModifyFilePath(filepath string, title string) error {
	filename := slug.Make(title) + ".md"
	newPath := path.Join("content", "posts", filename)
	return os.Rename(filepath, newPath)
}
