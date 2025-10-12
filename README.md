# Hugotui - Terminal User Interface for Hugo CLI

Hugotui is a terminal user interface for gohugo cli tool for help with content management. Currently it is only distributed here on GitHub.

## Installation

First, install gohugo if you haven't already. You can find installation instructions on the [Hugo website](https://gohugo.io/getting-started/installing/).
Then, download the latest release of Hugotui. Currently, you have to build it from source, so be sure you have golang installed. 
Please refer golang installation instructions on the [Golang website](https://golang.org/doc/install).
You can do this by cloning the repository into your Hugo site folder, and running the following commands:

``` bash
go build
```

## Usage
Run command:

``` bash 
./hugotui
```

## Configuration
There is couple of configurations options available in `hugo.toml` file. You can set the following options:

``` toml
hugotuiRemoteDir = "/var/www/html" # This is default, set in the source code, but you can override it here.
```

Hugotui also uses `baseURL` from your `hugo.toml` file to determine the remote server address for deployment.
Additionally, you have to set up SSH keys for passwordless authentication with your remote server. Also set `.ssh/config` file for your remote server. Example:

``` config
Host myserver
    HostName example.com
    User myuser
    Port 22
```

This is required for the deploy function to work properly.
By default, hugotui uses EDITOR environment variable to open the editor. You can set it to your preferred editor, for example:

``` bash
export EDITOR=nvim
# or
export EDITOR=emacs

```

## Help

Key bindings:

- j/k - Navigate up/down
- Enter - Select
- q - Quit
- n - New content and open in editor
- e - Edit content
- o - Open content in editor
- p - Preview site. PLEASE NOTE: This will run `hugo server` command in the background, and you have to manually kill it.
- P - Publish site over SCP. Builds Hugo site and copies the public folder to the remote server using SCP.
