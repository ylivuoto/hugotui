package utils

import (
	"fmt"
	"os"
	"os/exec"
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
	if editor == "" {
		editor = "nvim"
	}
	cmd := exec.Command(editor, filePath)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// m.status = fmt.Sprintf("Opening %s...", filePath)
	if err := cmd.Run(); err != nil {
		fmt.Printf("Failed: %v", err)
	}
	// else {
	// 	m.status = fmt.Sprintf("Closed editor for %s", filePath)
	// }
	// TODO: return error properly
	return nil
}
