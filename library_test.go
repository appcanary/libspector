package libspector

import (
	"bytes"
	"fmt"
	"testing"
)

const sampleDpkgOutput = `libept1.4.12: /usr/lib/libept.so.1.0.5.4.12
libpython2.7:amd64: /usr/lib/x86_64-linux-gnu/libpython2.7.so.1
libc6:amd64: /lib/x86_64-linux-gnu/libthread_db.so.1
libnss3-1d:amd64: /usr/lib/x86_64-linux-gnu/libnssutil3.so.1d
liblxc0: /usr/lib/x86_64-linux-gnu/liblxc.so.1.0.0.alpha1
`

func TestParseFindLibrary(t *testing.T) {
	var buf = new(bytes.Buffer)
	fmt.Fprint(buf, sampleDpkgOutput)

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
		t.Errorf("parsed wrong number of libraries; got %q; want %q", len(libs), len(cases))
	}

	for i, expected := range cases {
		lib := libs[i]
		if got := lib.Path(); got != expected.Path {
			t.Errorf("case #%d Path: got %q; want %q", i, got, expected.Path)
		}
	}
}
