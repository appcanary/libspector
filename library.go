package libspector

import (
	"bufio"
	"bytes"
	"errors"
	"os"
	"os/exec"
	"strings"
	"time"
)

var ErrParseCmdOutput = errors.New("failed to parse command output")

type library struct {
	path     string
	pkg      string
	modified *time.Time
}

// Path returns the full path of this library file.
func (lib *library) Path() string {
	return lib.path
}

// Modified returns the modification time of the library path.
func (lib *library) Modified() (time.Time, error) {
	if lib.modified != nil {
		return *lib.modified, nil
	}
	info, err := os.Stat(lib.path)
	if err != nil {
		return time.Time{}, err
	}
	mtime := info.ModTime()
	lib.modified = &mtime
	return mtime, nil
}

// parseFindLibrary parses output produced by commands run within FindLibrary,
// separated out for testing.
func parseFindLibrary(buf *bytes.Buffer) ([]library, error) {
	libs := []library{}
	scanner := bufio.NewScanner(buf)
	for scanner.Scan() {
		// Each line should look like:
		// somelib-1.0:maybeplatform: /usr/lib/somelib.so
		line := scanner.Text()
		parts := strings.Split(line, ":")
		if len(parts) < 2 {
			return nil, ErrParseCmdOutput
		}
		pkg, path := parts[0], strings.TrimSpace(parts[len(parts)-1])
		libs = append(libs, library{path: path, pkg: pkg})
	}
	if err := scanner.Err(); err != nil {
		return libs, err
	}

	return libs, nil
}

// FindLibrary uses `dpkg -S` to find libraries with the given path substring.
func FindLibrary(path string) ([]library, error) {
	cmd := exec.Command("dpkg", "-S", path)
	buf := new(bytes.Buffer)
	cmd.Stdout = buf

	if err := cmd.Run(); err != nil {
		return nil, err
	}

	return parseFindLibrary(buf)
}
