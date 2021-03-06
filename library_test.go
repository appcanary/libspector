package libspector

import (
	"bytes"
	"fmt"
	"testing"
	"time"
)

var sampleOutputs = map[string]string{
	"dpkg -S foo": `libept1.4.12: /usr/lib/libept.so.1.0.5.4.12
libpython2.7:amd64: /usr/lib/x86_64-linux-gnu/libpython2.7.so.1
libc6:amd64: /lib/x86_64-linux-gnu/libthread_db.so.1
libnss3-1d:amd64: /usr/lib/x86_64-linux-gnu/libnssutil3.so.1d
liblxc0: /usr/lib/x86_64-linux-gnu/liblxc.so.1.0.0.alpha1
`,
	"pmap 1234": `1234:   nginx: master process /usr/sbin/nginx
0000000000400000    788K r-x--  /usr/sbin/nginx
00000000006c4000      4K r----  /usr/sbin/nginx
00000000006c5000     84K rw---  /usr/sbin/nginx
00000000006da000     60K rw---    [ anon ]
0000000000cd3000    504K rw---    [ anon ]
0000000000d51000    844K rw---    [ anon ]
00007fee20e23000     48K r-x--  /lib/x86_64-linux-gnu/libnss_files-2.17.so
00007fee20e2f000   2044K -----  /lib/x86_64-linux-gnu/libnss_files-2.17.so
00007fee2102e000      4K r----  /lib/x86_64-linux-gnu/libnss_files-2.17.so
00007fee2102f000      4K rw---  /lib/x86_64-linux-gnu/libnss_files-2.17.so
`,
}

func TestParseFindLibrary(t *testing.T) {
	var buf = new(bytes.Buffer)
	fmt.Fprint(buf, sampleOutputs["dpkg -S foo"])

	libs, err := parseFindLibrary(buf)
	if err != nil {
		t.Error(err)
	}
	cases := []struct {
		Path string
	}{
		{"/usr/lib/libept.so.1.0.5.4.12"},
		{"/usr/lib/x86_64-linux-gnu/libpython2.7.so.1"},
		{"/lib/x86_64-linux-gnu/libthread_db.so.1"},
		{"/usr/lib/x86_64-linux-gnu/libnssutil3.so.1d"},
		{"/usr/lib/x86_64-linux-gnu/liblxc.so.1.0.0.alpha1"},
	}
	if len(libs) != len(cases) {
		t.Errorf("parsed wrong number of libraries; got %d; want %d", len(libs), len(cases))
	}

	for i, expected := range cases {
		lib := libs[i]
		if got := lib.Path(); got != expected.Path {
			t.Errorf("case #%d Path: got %q; want %q", i, got, expected.Path)
		}
	}
}

func TestParseFindLibraryByPID(t *testing.T) {
	var buf = new(bytes.Buffer)
	fmt.Fprint(buf, sampleOutputs["pmap 1234"])

	libs, err := parseFindLibraryByPID(buf)
	if err != nil {
		t.Error(err)
	}
	cases := []struct {
		Path string
	}{
		{"/usr/sbin/nginx"},
		{"/lib/x86_64-linux-gnu/libnss_files-2.17.so"},
	}
	if len(libs) != len(cases) {
		t.Errorf("parsed wrong number of libraries; got %d; want %d", len(libs), len(cases))
	}

	for i, expected := range cases {
		lib := libs[i]
		if got := lib.Path(); got != expected.Path {
			t.Errorf("case #%d Path: got %q; want %q", i, got, expected.Path)
		}
	}
}

func TestLibraryOutdated(t *testing.T) {
	ctime := time.Date(2016, time.January, 1, 12, 0, 0, 0, time.UTC)
	var lib Library = &library{
		path:    "/usr/lib/foo.so",
		pkgName: "libfoo",
		ctime:   &ctime,
	}

	// Process is newer than library
	stime := time.Date(2016, time.January, 1, 13, 0, 0, 0, time.UTC)
	proc := &process{
		pid:     1234,
		started: &stime,
	}
	if actual, expected := lib.Outdated(proc), false; actual != expected {
		t.Errorf("got %t; want %t", actual, expected)
	}

	// Process is older than library
	*proc.started = time.Date(2016, time.January, 1, 11, 0, 0, 0, time.UTC)
	if actual, expected := lib.Outdated(proc), true; actual != expected {
		t.Errorf("got %v; want %v", actual, expected)
	}
}
