package utils

import (
	"github.com/pelletier/go-toml"
)

var (
	HugoRemote     string
	HugoRemotePort string
)

func SetupConfig() {
	tree, err := toml.LoadFile("hugo.toml")
	if err != nil {
		panic(err)
	}
	HugoRemote = tree.Get("hugotuiScpRemote").(string)
	HugoRemotePort = tree.Get("hugotuiScpPort").(string)
}
