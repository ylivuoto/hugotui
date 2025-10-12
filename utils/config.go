package utils

import (
	"fmt"
	"net/url"

	"github.com/pelletier/go-toml"
)

var (
	HugoRemote    string
	HugoRemoteDir string = "/var/www/html/"
)

func SetupConfig() {
	// Use hugo.toml for parsing config data
	tree, err := toml.LoadFile("hugo.toml")
	if err != nil {
		panic(err)
	}

	// Parse remote host, ssh config should be set respectively
	host := tree.Get("baseURL")
	if host == nil {
		fmt.Println("No baseURL found in hugo.toml")
		return
	}

	url, err := url.Parse(host.(string))
	if err != nil {
		fmt.Println("Error parsing HugoRemote URL:", err)
	}

	HugoRemote = url.Hostname()
	dir := tree.Get("hugotuiPublishDir")
	if dir != nil {
		HugoRemoteDir = dir.(string)
	}
}
