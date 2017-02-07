package libspector

import (
	"syscall"
	"time"
)

// Ctime returns the changed time of the library path.

func (lib *library) Ctime() (time.Time, error) {
	if lib.ctime != nil {
		return *lib.ctime, nil
	}

	var stat syscall.Stat_t
	err := syscall.Stat(lib.path, &stat)

	if err != nil {
		return time.Time{}, err
	}

	sec, nsec := stat.Ctimespec.Unix()
	ctime := time.Unix(sec, nsec)
	lib.ctime = &ctime

	return ctime, nil
}
