package yespower

/*
#cgo CFLAGS: -std=gnu99 -I${SRCDIR}/..
#include <yespower-c/sha256.c>
#include <yespower-c/yespower-opt.c>
#include <yespower-c/yespower.c>
*/
import "C"

import (
	"unsafe"
)

// Yespower algorithm version numbers, matching yespower_version_t in
// yespower.h. Named with a Version suffix to avoid colliding with the
// string-valued YESPOWER_0_5 / YESPOWER_1_0 constants in yespower_native.go.
const (
	Version05 uint32 = 5
	Version10 uint32 = 10
)

// HashVersion computes the 32-byte yespower hash of input using the given
// version (YESPOWER_0_5 or YESPOWER_1_0), N, r and personalization token.
//
// It avoids the per-call C heap allocations of the previous implementation:
//   - input and per are passed by pointer into the Go-managed buffers, which
//     the C code only reads during the call (no copy into the C heap);
//   - the 32-byte digest is written directly into a Go-allocated array, so no
//     C malloc/free/GoBytes round-trip is needed.
//
// The Go memory passed here contains no Go pointers and is not retained by C
// after yespower_hash returns, so this is safe under the cgo pointer rules.
func HashVersion(input []byte, N uint32, r uint32, per string, version uint32) []byte {
	// Single allocation: the digest is written straight into this slice by
	// the C side, removing the previous malloc + GoBytes copy.
	hashed := make([]byte, 32)

	// cgo forbids taking &slice[0] / pointer into the string data when the
	// length is zero, so route empty buffers through a NULL pointer. The C
	// side already treats a zero length / NULL pers as "no personalization".
	var inPtr *C.char
	if len(input) > 0 {
		inPtr = (*C.char)(unsafe.Pointer(&input[0]))
	}

	var perPtr *C.char
	if len(per) > 0 {
		perPtr = (*C.char)(unsafe.Pointer(unsafe.StringData(per)))
	}

	C.yespower_hash(
		inPtr, C.uint(len(input)),
		C.uint(N), C.uint(r),
		perPtr, C.uint(len(per)),
		C.uint(version),
		(*C.char)(unsafe.Pointer(&hashed[0])),
	)

	return hashed
}

// Hash computes the YESPOWER_1_0 hash, preserving the original API.
func Hash(input []byte, N uint32, r uint32, per string) []byte {
	return HashVersion(input, N, r, per, Version10)
}
