package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/josh-baltar/nrebuild/internal/rebuild"
)

// manifest is the on-disk shape of a targets manifest.
type manifest struct {
	Targets []rebuild.Target `json:"targets"`
}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: nrebuild <manifest.json> [changed-path...]")
		os.Exit(2)
	}
	raw, err := os.ReadFile(os.Args[1])
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	var m manifest
	if err := json.Unmarshal(raw, &m); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	for _, name := range rebuild.SelectRebuilds(m.Targets, os.Args[2:]) {
		fmt.Println(name)
	}
}
