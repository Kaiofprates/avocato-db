package integrity

import (
	"crypto/sha256"
	"fmt"
)

// MMR represents a Merkle Mountain Range
type MMR struct {
	peaks [][]byte
	count uint64
}

func NewMMR() *MMR {
	return &MMR{
		peaks: make([][]byte, 0),
		count: 0,
	}
}

func (m *MMR) Append(hash []byte) {
	m.count++
	m.peaks = append(m.peaks, hash)
	m.mergePeaks()
}

func (m *MMR) mergePeaks() {
	// A simple MMR implementation:
	// We represent the tree structure using the binary representation of the count.
	// Each set bit in count corresponds to a peak.
	// We merge from right to left while we have peaks of the same height.
	// For simplicity in this POC, we'll keep the peaks list and merge the last two if they are at the same height.
	// Actually, a simpler way is to maintain the peaks based on the bits of the count.
	
	// If we just appended, we added a peak of height 0.
	height := 0
	
	for len(m.peaks) >= 2 {
		right := m.peaks[len(m.peaks)-1]
		left := m.peaks[len(m.peaks)-2]
		
		// If the number of leaves before this append was such that we need to merge
		// We can determine if we need to merge by checking if the bit at 'height' in (count-1) is 1.
		// If it's 1, it means there is an existing peak of this height.
		if ((m.count - 1) & (1 << height)) != 0 {
			// Merge
			m.peaks = m.peaks[:len(m.peaks)-2] // Remove last two
			
			merged := sha256.Sum256(append(left, right...))
			m.peaks = append(m.peaks, merged[:])
			height++
		} else {
			break
		}
	}
}

func (m *MMR) Root() []byte {
	if len(m.peaks) == 0 {
		return nil
	}
	
	// Bag the peaks
	root := m.peaks[len(m.peaks)-1]
	for i := len(m.peaks) - 2; i >= 0; i-- {
		merged := sha256.Sum256(append(m.peaks[i], root...))
		root = merged[:]
	}
	
	return root
}

func (m *MMR) RootHex() string {
	root := m.Root()
	if root == nil {
		return ""
	}
	return fmt.Sprintf("%x", root)
}
