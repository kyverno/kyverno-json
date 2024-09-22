package matching

import (
	"sync"

	"github.com/cespare/xxhash/v2"
	"github.com/elastic/go-freelru"
	"github.com/kyverno/kyverno-json/pkg/core/assertion"
	"github.com/kyverno/kyverno-json/pkg/core/compilers"
)

type _compilers = compilers.Compilers

type Compiler struct {
	_compilers
	*freelru.SyncedLRU[string, func() (assertion.Assertion, error)]
}

func hashStringXXHASH(s string) uint32 {
	sum := xxhash.Sum64String(s)
	return uint32(sum) //nolint:gosec
}

func NewCompiler(compiler compilers.Compilers, cacheSize uint32) Compiler {
	out := Compiler{
		_compilers: compiler,
	}
	if cache, err := freelru.NewSynced[string, func() (assertion.Assertion, error)](cacheSize, hashStringXXHASH); err == nil {
		out.SyncedLRU = cache
	}
	return out
}

func (c Compiler) CompileAssertion(hash string, value any, defaultCompiler string) (assertion.Assertion, error) {
	if c.SyncedLRU == nil || hash == "" {
		return assertion.Parse(value, c._compilers, defaultCompiler)
	}
	entry, _ := c.SyncedLRU.Get(hash)
	if entry == nil {
		entry = sync.OnceValues(func() (assertion.Assertion, error) {
			return assertion.Parse(value, c._compilers, defaultCompiler)
		})
		c.SyncedLRU.Add(hash, entry)
	}
	return entry()
}
