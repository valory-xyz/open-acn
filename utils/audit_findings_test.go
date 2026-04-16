// Tests in this file each correspond to a Critical or High finding in
// audits/AUDIT-2026-04-15.md. Every test is expected to FAIL until the
// underlying bug is fixed.

package utils

import (
	"encoding/hex"
	"testing"
)

// TestAuditC3_RecoverEthereumSignatureShortBytesPanics verifies that
// RecoverAddressFromEthereumSignature returns an error rather than
// panicking when given an undersized signature. See AUDIT-2026-04-15.md C3.
//
// The current implementation indexes sigBytes[64] without a length
// check; any signature shorter than 65 bytes triggers an out-of-range
// panic that any unauthenticated peer can deliver.
func TestAuditC3_RecoverEthereumSignatureShortBytesPanics(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("AUDIT C3: RecoverAddressFromEthereumSignature panicked on short signature input: %v", r)
		}
	}()

	// 30 bytes — below the required 65, so sigBytes[64] is OOB.
	short := "0x" + hex.EncodeToString(make([]byte, 30))
	_, err := RecoverAddressFromEthereumSignature([]byte("anything"), short)
	if err == nil {
		t.Fatalf("AUDIT C3: expected error on short signature, got nil")
	}
}

// TestAuditC3b_EthereumAddressFromPublicKeyShortInputPanics verifies
// that EthereumAddressFromPublicKey rejects an undersized public key
// input rather than panicking on `publicKey[2:]`. See AUDIT-2026-04-15.md C3
// (companion finding noted in the same section).
func TestAuditC3b_EthereumAddressFromPublicKeyShortInputPanics(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("AUDIT C3: EthereumAddressFromPublicKey panicked on short input: %v", r)
		}
	}()

	_, err := EthereumAddressFromPublicKey("0")
	if err == nil {
		t.Fatalf("AUDIT C3: expected error on short public key, got nil")
	}
}
