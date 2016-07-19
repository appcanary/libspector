package libspector

import (
	"bytes"
	"fmt"
	"io"
	"os/exec"
	"time"
)

type process struct {
	pid     int
	started *time.Time
}

func (p *process) PID() int {
	return p.pid
}

func FindProcess(command string) ([]process, error) {
	cmd := exec.Command("pgrep", command)
	buf := new(bytes.Buffer)
	cmd.Stdout = buf

	if err := cmd.Run(); err != nil {
		return nil, err
	}

	procs := []process{}
	for {
		var pid int
		_, err := fmt.Fscanf(buf, "%d\n", &pid)
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		procs = append(procs, process{pid: pid})
	}
	return procs, nil
}
