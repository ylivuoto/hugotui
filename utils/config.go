package utils

import "os"

var (
	HugoProject string
	HugoRemote  string
)

func SetupConfig() {
	HugoProject = os.Getenv("HUGO_PROJECT")
	HugoRemote = os.Getenv("HUGO_REMOTE")

	if HugoProject == "" {
		HugoProject = "."
	}
}
