package nixtui

import (
	"os"
	"path/filepath"
)

var homeDir = os.Getenv("HOME")

var nixpkgs string

func init() {
	nixpkgs = os.Getenv("NIXPKGS")
	if nixpkgs == "" {
		nixpkgs = filepath.Join(homeDir, ".nix-defexpr/channels/nixpkgs")
	}
}
