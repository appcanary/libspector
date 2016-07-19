package libspector

import (
	"testing"
	"time"
)

func TestParseLStart(t *testing.T) {
	_, err := parseLStart("Sun 17 Jul 13:30:26 2016")
	if err != nil {
		t.Error(err)
	}
}

func TestFindProcess(t *testing.T) {
	// Find some shells, any shells. Hopefully there are shells.
	query := "sh"

	procs, err := FindProcess(query)
	if err != nil {
		t.Error(err)
	}

	if len(procs) == 0 {
		t.Errorf("failed to find process")
	}

	noTime := time.Time{}
	for i, proc := range procs {
		pid := proc.PID()
		if pid == 0 {
			t.Errorf("invalid pid at %d of %d", i, len(procs))
		}
		started, err := proc.Started()
		if err != nil {
			t.Errorf("err getting start time for pid %d: %s", pid, err)
		}
		if started == noTime {
			t.Errorf("missing start time for pid: %d", pid)
		}
	}
}
