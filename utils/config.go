package utils

import "os"

var HugoRemote string

func SetupConfig() {
	HugoRemote = os.Getenv("HUGO_REMOTE")
}
