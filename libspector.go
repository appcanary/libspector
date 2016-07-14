package libspector

import "time"

type Package interface {
	Name() string
	Version() string
}

type Library interface {
	Path() string
	Modified() time.Time
	Outdated() bool

	// Distribution package manager's dependency that owns this library.
	Package() Package

	// Find processes using this library
	Processes() []Process
}

type Process interface {
	PID() int
	Started() time.Time

	// Find libraries used by this process
	Libraries() []Library
}

type Query interface {
	FindProcess(command string) ([]Process, error)
	FindLibrary(path string) ([]Library, error)
}
