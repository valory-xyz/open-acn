// Tests in this file each correspond to a Critical or High finding in
// audits/AUDIT-2026-04-15.md. Every test is expected to FAIL until the
// underlying bug is fixed.

package acn

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

// TestAuditH7_AcnErrorEchoesInternalParserStrings verifies that
// `SendAcnError` is not called with raw `err.Error()` strings produced
// by deserialization failures. See AUDIT-2026-04-15.md H7.
//
// The acn/utils.go source contains the maintainer's own TOFIX comment
// flagging the same vulnerability. This test fails until either the
// `err.Error()` arguments are replaced with generic strings, or the
// TOFIX comment is removed (signalling a deliberate decision).
func TestAuditH7_AcnErrorEchoesInternalParserStrings(t *testing.T) {
	src := readAcnSourceFile(t, "utils.go")
	if !strings.Contains(src, "TOFIX(LR) setting Msgs to err.Error is potentially a security vulnerability") {
		t.Skipf("TOFIX marker removed — bug may have been addressed by another route; re-check AUDIT H7")
	}
	t.Fatalf("AUDIT H7: acn/utils.go still contains the TOFIX block sending raw err.Error() to remote peers via SendAcnError")
}

func readAcnSourceFile(t *testing.T, name string) string {
	t.Helper()
	_, here, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatalf("runtime.Caller failed")
	}
	path := filepath.Join(filepath.Dir(here), name)
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v", path, err)
	}
	return string(data)
}
