package utils

import (
	"encoding/binary"
	"math/rand"
	"strings"
	"time"

	cryptoRand "crypto/rand"
)

const (
	AlphaNumericSymbol       string = "abcdefghijklmnopqrstuvwxyz" + "ABCDEFGHIJKLMNOPRQSTUVWXYZ" + "0123456789" + "-_"
	alphaNumericSymbolLength int    = len(AlphaNumericSymbol)
)

// Seeds the default math.rand generator with a cryptographically secure random value
// Returns an error in cases where a system cryptographically secure PRNG is not available
func SeedRandom() error {
	// Prepare Seed
	// > Time-based seeds
	var seed int64 = time.Now().UnixMicro() ^ time.Now().UnixNano()
	// > System cryptographically secure PRNG
	b := make([]byte, 8)
	_, err := cryptoRand.Read(b)
	if err != nil {
		// Seed with time alone (weak) and return an error
		rand.Seed(seed)
		return err
	}
	seed ^= int64(binary.LittleEndian.Uint64(b))
	// Seed
	rand.Seed(seed)
	return nil
}

// Returns a slice of random bytes of a given length
func slowRandomBytes(b []byte) []byte {
	for i := len(b) - 1; i >= 0; i-- {
		b[i] = byte(rand.Intn(256))
	}
	return b
}

// Returns a slice of random bytes of a given length
// First attempts to use a reader, on error will fallback on slowRandomBytes
func RandomBytes(length int) []byte {
	b := make([]byte, length)
	_, err := rand.Read(b)
	if err != nil {
		return slowRandomBytes(b)
	}
	return b
}

// Returns a random string of a given length
// Uses AlphaNumericSymbol alphabet
func RandomString(length int) string {
	srcBytes := RandomBytes(length)
	builder := strings.Builder{}
	builder.Grow(length)
	for i := 0; i < length; i++ {
		builder.WriteByte(AlphaNumericSymbol[int(srcBytes[i])%alphaNumericSymbolLength])
	}
	return builder.String()
}

// Returns a slice of length [minLen, maxLen] with unique int in range [0, max)
//
// This must be true: (maxLen-minLen) <= max
func RandomIntSlice(max int, minLen int, maxLen int) []int {
	if (maxLen - minLen) > max {
		return []int{}
	}
	slice := make([]int, max)
	for i := 0; i < max; i++ {
		slice[i] = i
	}
	length := minLen
	if maxLen > minLen {
		length = minLen + rand.Intn(maxLen-minLen+1)
	}
	for i := 0; i < max-length; i++ {
		el := rand.Intn(len(slice))
		slice[el] = slice[len(slice)-1]
		slice = slice[:len(slice)-1]
	}
	return slice
}
