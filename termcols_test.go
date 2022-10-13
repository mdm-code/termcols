package termcols

import (
	"testing"
)

func BenchmarkColorize(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Colorize("foo, bar, baz", BlackFg)
	}
}
