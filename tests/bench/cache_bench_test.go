package bench

import (
	"avocato-db/src/core/ledger"
	"fmt"
	"testing"
)

func BenchmarkCache_Get(b *testing.B) {
	cache := ledger.NewCache(1000)
	block := &ledger.Block{Index: 1}
	cache.Put("1", block)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = cache.Get("1")
	}
}

func BenchmarkCache_Put(b *testing.B) {
	cache := ledger.NewCache(1000)
	block := &ledger.Block{Index: 1}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cache.Put(fmt.Sprintf("%d", i), block)
	}
}
