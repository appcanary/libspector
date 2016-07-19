package libspector

import "testing"

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

	for i, proc := range procs {
		pid := proc.PID()
		if pid == 0 {
			t.Errorf("invalid pid at %d of %d", i, len(procs))
		}
	}
}
