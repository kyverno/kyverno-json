package matching

import (
	"sync"

	"github.com/cespare/xxhash"
	"github.com/elastic/go-freelru"
	"github.com/kyverno/kyverno-json/pkg/core/assertion"
	"github.com/kyverno/kyverno-json/pkg/core/templating"
)

type Compiler struct {
	templating.Compiler
	*freelru.SyncedLRU[string, func() (assertion.Assertion, error)]
}

func hashStringXXHASH(s string) uint32 {
	sum := xxhash.Sum64String(s)
	return uint32(sum) //nolint:gosec
}

func NewCompiler(compiler templating.Compiler, cacheSize uint32) Compiler {
	out := Compiler{
		Compiler: compiler,
	}
	if cache, err := freelru.NewSynced[string, func() (assertion.Assertion, error)](cacheSize, hashStringXXHASH); err == nil {
		out.SyncedLRU = cache
	}
	return out
}

func (c Compiler) CompileAssertion(hash string, value any) (assertion.Assertion, error) {
	if c.SyncedLRU == nil {
		return assertion.Parse(value, c.Compiler)
	}
	entry, _ := c.SyncedLRU.Get(hash)
	if entry == nil {
		entry = sync.OnceValues(func() (assertion.Assertion, error) {
			return assertion.Parse(value, c.Compiler)
		})
		c.SyncedLRU.Add(hash, entry)
	}
	return entry()
}
