package libspector

import (
	"syscall"
	"time"
)

// Modified returns the modification time of the library path.

func (lib *library) Modified() (time.Time, error) {
	if lib.modified != nil {
		return *lib.modified, nil
	}

	var stat syscall.Stat_t
	err := syscall.Stat(lib.path, &stat)

	if err != nil {
		return time.Time{}, err
	}

	sec, nsec := stat.Ctimespec.Unix()
	ctime := time.Unix(sec, nsec)
	lib.modified = &ctime

	return ctime, nil
}
