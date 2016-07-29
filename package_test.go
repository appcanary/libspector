package libspector

import (
	"bytes"
	"fmt"
	"testing"
)

const sampleAptCachePolicyOutput = `libc6:
  Installed: 2.17-93ubuntu4
  Candidate: 2.17-93ubuntu4
  Version table:
 *** 2.17-93ubuntu4 0
        500 http://archive.ubuntu.com/ubuntu/ saucy/main amd64 Packages
        100 /var/lib/dpkg/status
`

func TestParseFindPackage(t *testing.T) {
	var buf = new(bytes.Buffer)
	fmt.Fprint(buf, sampleAptCachePolicyOutput)

	pkg, err := parseFindPackage(buf)
	if err != nil {
		t.Fatal(err)
	}
	if want, got := "libc6", pkg.Name(); want != got {
		t.Errorf("want %q; got %q", want, got)
	}
	if want, got := "2.17-93ubuntu4", pkg.Version(); want != got {
		t.Errorf("want %q; got %q", want, got)
	}
}
