package libspector

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"os/exec"
	"strings"
)

var ErrNoPackage = errors.New("no package provided")

// pkg represents a package installed by the distribution's package manager.
type pkg struct {
	name    string
	version string
}

// Name of the package
func (p *pkg) Name() string {
	return p.name
}

// Version of the package
func (p *pkg) Version() string {
	return p.version
}

func parseFindPackage(buf *bytes.Buffer) (*pkg, error) {
	scanner := bufio.NewScanner(buf)
	scanner.Scan()
	name := strings.TrimRight(scanner.Text(), ":")
	if name == "" {
		return nil, ErrParseCmdOutput
	}

	scanner.Scan()
	parts := strings.Split(scanner.Text(), ": ")
	if len(parts) != 2 {
		return nil, ErrParseCmdOutput
	}
	version := parts[1]

	return &pkg{name: name, version: version}, nil
}

// findPackage uses `apt-cache policy $name` to load package info
func findPackage(name string) (*pkg, error) {
	if name == "" {
		return nil, ErrNoPackage
	}
	cmd := exec.Command("apt-cache", "policy", name)
	buf := new(bytes.Buffer)
	cmd.Stdout = buf

	if err := cmd.Run(); err != nil {
		return nil, err
	}

	pkg, err := parseFindPackage(buf)
	if err != nil {
		return nil, err
	}
	if pkg.Name() != name {
		return pkg, fmt.Errorf("parsed package mismatch: %q != %q", pkg.Name(), name)
	}
	return pkg, nil
}
