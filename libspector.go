// Package libspector provides query functions for finding installed libraries
// and libraries used by active processes.
//
// This functionality is implemented by running various platform-native
// commands, concurrency safety is not implemented.
package libspector

import "time"

// Package that is managed by the distribution's package manager.
type Package interface {
	// Name of the package.
	Name() string

	// Version of the package installed.
	Version() string
}

// Library is a file representing a dynamically linked library or shared object.
type Library interface {
	// Path returns the absolute path of the library on the filesystem.
	Path() string

	// Ctime returns the last ctime of the library on the filesystem.
	Ctime() (time.Time, error)

	// Outdated returns whether Process was started earlier than the Ctime time of this library.
	Outdated(Process) bool

	// Distribution package manager's dependency that owns this library.
	Package() (Package, error)
}

// Process is a currently-running process.
type Process interface {
	// PID returns the process ID.
	PID() int

	// Started returns the time when the process was started, if it's still running.
	Started() (time.Time, error)

	// Find libraries used by this process
	Libraries() ([]Library, error)

	// Command line used to start the process
	CommandArgs() (string, error)

	CommandName() (string, error)
}

type Query interface {
	// AllProcesses returns all the running processes on the system.
	AllProcesses() ([]Process, error)

	// FindProcess finds all running processes that match the command substring.
	FindProcess(command string) ([]Process, error)

	// FindLibrary finds all the installed libraries that match the path substring.
	FindLibrary(path string) ([]Library, error)
}
