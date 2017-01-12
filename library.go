package libspector

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"
)

type library struct {
	path     string
	pkgName  string
	pkg      Package
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
	if lib.pkg != nil {
		// Is the package already loaded?
		return lib.pkg, nil
	}
	if lib.pkgName != "" {
		// Do we have the package name?
		pkg, err := findPackage(lib.pkgName)
		if err != nil {
			return nil, err
		}
		lib.pkg = pkg
		return pkg, nil
	}
	// Find the package from scratch.
	libs, err := FindLibrary(lib.path)
	if err != nil {
		return nil, err
	}
	if len(libs) != 1 {
		return nil, fmt.Errorf("failed to uniquely identify path %q, found %d results", lib.path, len(libs))
	}
	return libs[0].Package()
}

// Outdated compares the modified time of the library path against the timestamp of when the process was started.
func (lib *library) Outdated(proc Process) bool {
	mtime, err := lib.Modified()
	if err != nil {
		// Library path could not be queried, must be outdated.
		return true
	}
	stime, err := proc.Started()
	if err != nil {
		// Process start time could not be queried, must have been killed so can't be outdated.
		return false
	}
	return stime.Before(mtime)
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
			return nil, ErrParse(line)
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
		if strings.HasPrefix(line, " ") {
			// Finished parsing, found " total: ..." line
			break
		}
		parts := strings.Fields(line)
		if len(parts) < 4 {
			return nil, ErrParse(line)
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

// findLibraryByPID uses `pmap -p $PID` to find libraries that are being used by a given PID.
func findLibraryByPID(pid int) ([]Library, error) {
	cmd := exec.Command("pmap", "-p", fmt.Sprintf("%d", pid))
	buf := new(bytes.Buffer)
	cmd.Stdout = buf

	if err := cmd.Run(); err != nil {
		return nil, err
	}

	return parseFindLibraryByPID(buf)
}
