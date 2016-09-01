package libspector

import (
	"bytes"
	"fmt"
	"io"
	"os/exec"
	"strings"
	"time"
)

func parseLStart(s string) (time.Time, error) {
	return time.Parse(layoutLStart, s)
}

func ProcessByPID(pid int) Process {
	return &process{
		pid: pid,
	}
}

type process struct {
	pid     int
	started *time.Time
	command string
}

// PID returns the process ID.
func (p *process) PID() int {
	return p.pid
}

// Started uses `ps` to query the start timestamp of the process.
func (p *process) Started() (time.Time, error) {
	if p.started != nil {
		return *p.started, nil
	}

	started := time.Time{}
	cmd := exec.Command("ps", "-p", fmt.Sprintf("%d", p.PID()), "-o", "lstart=")
	out, err := cmd.Output()
	if err != nil {
		return started, err
	}

	started, err = parseLStart(strings.TrimSpace(string(out)))
	if err != nil {
		return started, err
	}
	p.started = &started
	return started, nil
}

// Command uses `ps` to query the command of the process.
func (p *process) Command() (string, error) {
	if p.command != "" {
		return p.command, nil
	}

	cmd := exec.Command("ps", "-p", fmt.Sprintf("%d", p.PID()), "-o", "cmd=")
	out, err := cmd.Output()
	if err != nil {
		return "", err
	}

	command := strings.TrimSpace(string(out))
	p.command = command
	return command, nil
}

// Libraries returns the dynamically linked libraries used by this process.
func (p *process) Libraries() ([]Library, error) {
	return findLibraryByPID(p.pid)
}

// FindProcess uses `pgrep` to find all processes that match a command.
func FindProcess(command string) ([]Process, error) {
	// TODO: Do we want more flexible querying abilities? Such as full arg substring, or parent pid, etc?
	// TODO: Consider getting the full process list at once, including start time with `ps aux` or similar?
	cmd := exec.Command("pgrep", command)
	buf := new(bytes.Buffer)
	cmd.Stdout = buf

	if err := cmd.Run(); err != nil {
		return nil, err
	}

	procs := []Process{}
	for {
		var pid int
		_, err := fmt.Fscanf(buf, "%d\n", &pid)
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		procs = append(procs, &process{pid: pid})
	}
	return procs, nil
}
