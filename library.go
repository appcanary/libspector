package libspector

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"
)

var ErrParseCmdOutput = errors.New("failed to parse command output")

type library struct {
	path     string
	pkgName  string
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

func (lib *library) Package() (Package, error) {
	return findPackage(lib.pkgName)
}

func (lib *library) Outdated() bool {
	// XXX: Implement stub
	return false
}

// parseFindLibrary parses output produced by commands run within FindLibrary,
// separated out for testing.
func parseFindLibrary(buf *bytes.Buffer) ([]Library, error) {
	libs := []Library{}
	scanner := bufio.NewScanner(buf)
	for scanner.Scan() {
		// Each line should look like:
		// somelib-1.0:maybearch: /usr/lib/somelib.so
		line := scanner.Text()
		parts := strings.Split(line, ":")
		if len(parts) < 2 {
			return nil, ErrParseCmdOutput
		}
		pkg, path := parts[0], strings.TrimSpace(parts[len(parts)-1])
		libs = append(libs, &library{path: path, pkgName: pkg})
	}
	if err := scanner.Err(); err != nil {
		return libs, err
	}

	return libs, nil
}

// FindLibrary uses `dpkg -S` to find libraries with the given path substring.
func FindLibrary(path string) ([]Library, error) {
	cmd := exec.Command("dpkg", "-S", path)
	buf := new(bytes.Buffer)
	cmd.Stdout = buf

	if err := cmd.Run(); err != nil {
		return nil, err
	}

	return parseFindLibrary(buf)
}

func parseFindLibraryByPID(buf *bytes.Buffer) ([]Library, error) {
	// First line looks like this:
	// 1234:    nginx: master process /usr/sbin/nginx
	// Then every following line is:
	// 0000000000400000    788K r-x--  /usr/sbin/nginx
	libs := []Library{}
	scanner := bufio.NewScanner(buf)
	scanner.Scan() // Skip first line

	seen := map[string]bool{}

	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Fields(line)
		if len(parts) < 4 {
			return nil, ErrParseCmdOutput
		}
		path := parts[3]
		if !strings.HasPrefix(path, "/") {
			// Skip anon lines etc.
			continue
		}
		if _, ok := seen[path]; ok {
			// Already seen, skip
			continue
		}
		seen[path] = true
		libs = append(libs, &library{path: path})
	}
	if err := scanner.Err(); err != nil {
		return libs, err
	}

	return libs, nil
}

// findLibraryByPID uses `pldd $PID` to find libraries that are being used by a given PID.
func findLibraryByPID(pid int) ([]Library, error) {
	cmd := exec.Command("pldd", fmt.Sprintf("%d", pid))
	buf := new(bytes.Buffer)
	cmd.Stdout = buf

	if err := cmd.Run(); err != nil {
		return nil, err
	}

	return parseFindLibrary(buf)
}
