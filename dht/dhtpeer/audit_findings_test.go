// Tests in this file each correspond to a Critical or High finding in
// audits/AUDIT-2026-04-15.md. Every test is expected to FAIL until the
// underlying bug is fixed.
//
// Several of these findings require a full DHTPeer + delegate + mailbox
// harness to reproduce end-to-end. Where that scaffolding does not yet
// exist in this repo, the test stops with a clear `AUDIT <ID>` failure
// pointing back to the audit, so the finding remains visible in CI.
//
// Run only these tests:
//   go test -gcflags=-l -count=1 -v -run TestAudit ./dht/dhtpeer/...

package dhtpeer

import (
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	"testing"
)

// ---------------------------------------------------------------------
// Static-pattern tests (decisive, fast)
// ---------------------------------------------------------------------

// TestAuditN7_FmtPrintlnSilentDiscardOnSyncQueueFull verifies that the
// "channel full, discarding" branch in the per-pair enqueue path does
// not silently drop envelopes via fmt.Println. See AUDIT-2026-04-15.md N7.
func TestAuditN7_FmtPrintlnSilentDiscardOnSyncQueueFull(t *testing.T) {
	src := readDhtpeerSource(t, "dhtpeer.go")
	if strings.Contains(src, `fmt.Println("CHANNEL FULL, DISCARDING`) {
		t.Fatalf("AUDIT N7: dhtpeer.go discards envelopes via fmt.Println without using the structured logger or any metric")
	}
}

// TestAuditN8_RouteEnvelopeDHTLookupUsesContextBackground verifies that
// the DHT lookup helper does not use `context.Background()` (which
// cannot be cancelled when the peer shuts down). See AUDIT-2026-04-15.md N8.
func TestAuditN8_RouteEnvelopeDHTLookupUsesContextBackground(t *testing.T) {
	src := readDhtpeerSource(t, "dhtpeer.go")
	// Look for a context.Background() call in _routeEnvelopeDHTLookup.
	// We grep with a +/- 30 line window around the function header.
	idx := strings.Index(src, "_routeEnvelopeDHTLookup")
	if idx < 0 {
		t.Skipf("function _routeEnvelopeDHTLookup not found — refactored?")
	}
	end := idx + 4000
	if end > len(src) {
		end = len(src)
	}
	if strings.Contains(src[idx:end], "context.Background()") {
		t.Fatalf("AUDIT N8: _routeEnvelopeDHTLookup constructs context.Background() instead of deriving from a parent that cancels on shutdown (dhtPeer.closing)")
	}
}

// TestAuditH6_PrivateKeyLoggedAtDebugLevel verifies that the bootstrap
// path does not log the AEA private key. See AUDIT-2026-04-15.md H6.
//
// The bug lives in aea/api.go but the example output is already
// documented in README.md, so a static check on the README is the
// cheapest way to keep this finding visible.
func TestAuditH6_PrivateKeyLoggedAtDebugLevel(t *testing.T) {
	apiSrc := readRepoFile(t, "aea/api.go")
	// The exact line is `logger.Debug().Msgf("id: %s", aea.id)`.
	if regexp.MustCompile(`logger\.Debug\(\)\.Msgf\(["']id: %s["']`).MatchString(apiSrc) {
		t.Fatalf("AUDIT H6: aea/api.go logs the AEA private key at debug level (`id: %%s`)")
	}
}

// ---------------------------------------------------------------------
// Skeleton tests for findings that require a network harness
// ---------------------------------------------------------------------

// TestAuditC4_MailboxSessionIdIsUnboundBearerToken — see AUDIT C4.
//
// To prove the bug requires: spinning up a MailboxServer, registering
// agent A and capturing its Session-Id, then issuing a GET
// /get_envelope from a *different* TLS client that supplies the
// captured Session-Id and asserting the request is rejected. Today
// it succeeds and drains A's mailbox.
func TestAuditC4_MailboxSessionIdIsUnboundBearerToken(t *testing.T) {
	t.Fatal("AUDIT C4: not yet reproduced in a runtime test — see audits/AUDIT-2026-04-15.md C4. Mailbox session tokens are bearer tokens with no binding to the registering client (TLS cert, IP, or PoR).")
}

// TestAuditH1_DelegateConnectionHasNoReadDeadline — see AUDIT H1.
//
// Reproduction sketch: bring up a DHTPeer with delegate enabled, dial
// the delegate port, send no bytes, sleep > 30s, assert the connection
// has been closed by the peer. Currently it stays open indefinitely
// (slowloris).
func TestAuditH1_DelegateConnectionHasNoReadDeadline(t *testing.T) {
	t.Fatal("AUDIT H1: not yet reproduced in a runtime test — see audits/AUDIT-2026-04-15.md H1. handleNewDelegationConnection sets no SetReadDeadline / SetWriteDeadline on accepted connections.")
}

// TestAuditH2_SyncMessagesMapLeaksAcrossConnectionLifetime — see AUDIT H2.
func TestAuditH2_SyncMessagesMapLeaksAcrossConnectionLifetime(t *testing.T) {
	t.Fatal("AUDIT H2: not yet reproduced in a runtime test — see audits/AUDIT-2026-04-15.md H2. handleNewDelegationConnection cleans tcpAddresses and acnStatuses on disconnect but never deletes from syncMessages or terminates the per-pair routing goroutines.")
}

// TestAuditH3_AcnStatusesSendIsBlocking — see AUDIT H3.
func TestAuditH3_AcnStatusesSendIsBlocking(t *testing.T) {
	t.Fatal("AUDIT H3: not yet reproduced in a runtime test — see audits/AUDIT-2026-04-15.md H3. AddAcnStatusMessage performs a blocking send to acnStatuses[addr]; a slow delegate client stalls the routing goroutine.")
}

// TestAuditH4_MailboxLockHeldAcrossNetworkIO — see AUDIT H4.
func TestAuditH4_MailboxLockHeldAcrossNetworkIO(t *testing.T) {
	t.Fatal("AUDIT H4: not yet reproduced in a runtime test — see audits/AUDIT-2026-04-15.md H4. apiGetEnvelope holds mailboxServer.lock across res.Write(buf), so a single slow HTTP reader blocks every mailbox operation for every agent.")
}

// TestAuditH5_AeaApiClosingFlagRace — see AUDIT H5.
//
// Should be exercised under `make race_test`. Even without -race, a
// concurrent double-close of send_queue / out_queue panics.
func TestAuditH5_AeaApiClosingFlagRace(t *testing.T) {
	t.Fatal("AUDIT H5: not yet reproduced in a runtime test — see audits/AUDIT-2026-04-15.md H5. aea/api.go uses a plain bool `closing` mutated from multiple goroutines and closes channels with no sync.Once.")
}

// TestAuditH8_PoRAcceptedAsStableBearerCredential — see AUDIT H8.
func TestAuditH8_PoRAcceptedAsStableBearerCredential(t *testing.T) {
	t.Fatal("AUDIT H8: not yet reproduced in a runtime test — see audits/AUDIT-2026-04-15.md H8. The registration handler accepts any AgentRecord whose signature verifies, with no nonce / replay / freshness check; a captured PoR can be replayed at any peer.")
}

// TestAuditH9_PersistentAgentRecordStorageGrowsUnbounded — see AUDIT H9.
func TestAuditH9_PersistentAgentRecordStorageGrowsUnbounded(t *testing.T) {
	t.Fatal("AUDIT H9: not yet reproduced in a runtime test — see audits/AUDIT-2026-04-15.md H9. saveAgentRecordToPersistentStorage appends without bound and initAgentRecordPersistentStorage reads the entire file into memory at startup.")
}

// TestAuditN6_PerPairOrderingViolatedOnDhtMiss — see AUDIT N6.
//
// The reviewer notes that TestMessageOrderingWithDelegateClientTwoHops
// was written to exercise this but currently passes because the
// libp2p v0.33 NullResourceManager change masks the slow-path entry.
// Reproducing the bug requires injecting an artificial DHT miss; the
// scaffolding is non-trivial.
func TestAuditN6_PerPairOrderingViolatedOnDhtMiss(t *testing.T) {
	t.Fatal("AUDIT N6: not yet reproduced in a runtime test — see audits/AUDIT-2026-04-15.md N6. RouteEnvelope returns immediately on DHT miss after pushing to slow_queue; per-pair goroutine moves to N+1 and a subsequent fast-path send can overtake N. README guarantees per-(sender, recipient) total ordering.")
}

// ---------------------------------------------------------------------
// helpers
// ---------------------------------------------------------------------

func readDhtpeerSource(t *testing.T, name string) string {
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

func readRepoFile(t *testing.T, rel string) string {
	t.Helper()
	_, here, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatalf("runtime.Caller failed")
	}
	// walk up to find the repo root (contains go.mod)
	dir := filepath.Dir(here)
	for i := 0; i < 6; i++ {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			path := filepath.Join(dir, rel)
			data, err := os.ReadFile(path)
			if err != nil {
				t.Fatalf("read %s: %v", path, err)
			}
			return string(data)
		}
		dir = filepath.Dir(dir)
	}
	t.Fatalf("could not locate repo root from %s", here)
	return ""
}
