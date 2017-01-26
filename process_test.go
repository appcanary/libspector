package libspector

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"runtime"
	"testing"
	"time"
)

func TestParseLStart(t *testing.T) {
	var err error
	switch runtime.GOOS {
	case "linux":
		_, err = parseLStart("Sun Jul 17 13:30:26 2016")
	case "darwin":
		_, err = parseLStart("Sun 17 Jul 13:30:26 2016")
	}

	if err != nil {
		t.Error(err)
	}
}

func TestAllProcesses(t *testing.T) {
	procs, err := AllProcesses()
	if err != nil {
		t.Error(err)
	}

	if len(procs) == 0 {
		t.Errorf("There should be at least an init process")
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

func TestProcessCommand(t *testing.T) {
	script, err := filepath.Abs("test_script.sh")
	if err != nil {
		t.Error(err)
	}

	cmd := exec.Command(script, "1", "hello", "douche", "4")
	err = cmd.Start()
	if err != nil {
		t.Error(err)
	}
	defer cmd.Process.Kill()

	procs, err := FindProcess("test_script")
	if err != nil {
		t.Error(err)
	}

	if len(procs) == 0 {
		t.Error("failed to find process")
	}

	fullCommand, err := procs[0].Command()
	if err != nil {
		t.Error(err)
	}

	if fullCommand != fmt.Sprintf("sh %s 1 hello douche 4", script) {
		t.Errorf("Command was wrong: %s", fullCommand)
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
		t.Error("failed to find process")
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
