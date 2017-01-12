package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"

	"github.com/appcanary/libspector"
)

func fail(code int, msg string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, msg, args...)
	os.Exit(code)
}

func main() {
	flag.Parse()

	pidArg := flag.Arg(0)
	if pidArg == "" {
		flag.Usage()
		os.Exit(2)
	}

	pid, err := strconv.Atoi(pidArg)
	if err != nil {
		fail(3, "PID must be an integer: %s\n", err)
	}

	proc := libspector.ProcessByPID(pid)
	started, err := proc.Started()
	if err != nil {
		fail(4, "PID %d is not running: %s\n", pid, err)
	}

	fmt.Printf("PID %d started %s\n", pid, started)

	libs, err := proc.Libraries()
	if err != nil {
		fail(5, "Failed to load libraries: %s\n", err)
	}

	fmt.Println("\nMemory Map:")
	for _, lib := range libs {
		if !lib.Outdated(proc) {
			fmt.Printf("  * %s\n", lib.Path())
			continue
		}

		pkg, err := lib.Package()
		if err != nil {
			fail(6, "Outdated lib %q failed to load packge: %s\n", lib.Path(), err)
		}
		fmt.Printf("! * %s [Outdated: %s %s]\n", lib.Path(), pkg.Name(), pkg.Version())
	}
}

func init() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "usage: %s PID\n", os.Args[0])
		flag.PrintDefaults()
	}
}
