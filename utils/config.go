package utils

import "os"

var HugoProject string

func SetupConfig() {
	HugoProject = os.Getenv("HUGO_PROJECT")

	if HugoProject == "" {
		HugoProject = "."
	}
}
