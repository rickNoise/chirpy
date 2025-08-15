package auth

import (
	"crypto/rand"
	"encoding/hex"
)

// Generates a random 256-bit (32-byte) hex-encoded string
func MakeRefreshToken() (string, error) {
	// generate 32 bytes (256 bits) of random data from the crypto/rand package
	randomSlice := make([]byte, 32)
	rand.Read(randomSlice)

	// convert the random data to a hex string
	return hex.EncodeToString(randomSlice), nil
}
