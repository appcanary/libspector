// Package libspector provides query functions for finding installed libraries
// and libraries used by active processes.
//
// This functionality is implemented by running various platform-native
// commands, concurrency safety is not implemented.
package libspector

import "time"

// Package that is managed by the distribution's package manager.
type Package interface {
	Name() string
	Version() string
}

// Library is a file representing a dynamically linked library or shared object.
type Library interface {
	Path() string
	Modified() (time.Time, error)
	Outdated() bool

	// Distribution package manager's dependency that owns this library.
	Package() (Package, error)
}

// Process is a currently-running process.
type Process interface {
	PID() int
	Started() (time.Time, error)

	// Find libraries used by this process
	Libraries() ([]Library, error)
}

type Query interface {
	FindProcess(command string) ([]Process, error)
	FindLibrary(path string) ([]Library, error)
}
