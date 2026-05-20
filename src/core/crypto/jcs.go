package crypto

import (
	"github.com/lattice-substrate/json-canon/jcs"
)

// Canonicalize returns the deterministic byte representation of the input according to RFC 8785.
func Canonicalize(data []byte) ([]byte, error) {
	return jcs.Canonicalize(data)
}
