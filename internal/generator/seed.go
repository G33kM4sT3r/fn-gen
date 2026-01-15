package generator

import (
	"crypto/sha256"
	"encoding/binary"
)

// HashToUint64 converts a string input into a deterministic 64-bit unsigned integer.
// It uses SHA-256 hashing to ensure uniform distribution across the output space.
//
// The function:
//  1. Computes the SHA-256 hash of the input string (256 bits / 32 bytes)
//  2. Takes the first 8 bytes of the hash
//  3. Interprets those bytes as a big-endian uint64
//
// This provides a deterministic mapping from any string to a number in the range
// [0, 2^64-1], which can then be used with modulo to select from word lists.
//
// Properties:
//   - Deterministic: same input always produces same output
//   - Well-distributed: SHA-256 ensures even distribution
//   - Collision-resistant: different inputs produce different outputs
func HashToUint64(input string) uint64 {
	// Compute SHA-256 hash (returns fixed 32-byte array)
	hash := sha256.Sum256([]byte(input))

	// Extract first 8 bytes as uint64 (big-endian byte order)
	return binary.BigEndian.Uint64(hash[:8])
}
