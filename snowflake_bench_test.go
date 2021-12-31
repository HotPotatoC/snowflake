package snowflake_test

import (
	"testing"

	"github.com/HotPotatoC/snowflake"
	bwmarrinsnowflake "github.com/bwmarrin/snowflake"
	gosnowflake "github.com/godruoyi/go-snowflake"
)

func BenchmarkNewID(b *testing.B) {
	benchmarks := []struct {
		name string
		fn   func(b *testing.B)
	}{
		{"github.com/HotPotatoC/snowflake", benchmarkHotPotatoCSnowflake},
		{"github.com/bwmarrin/snowflake", benchmarkBwmarrinSnowflake},
		{"github.com/godruoyi/go-snowflake", benchmarkGoSnowflake},
	}

	for _, bb := range benchmarks {
		b.Run(bb.name, bb.fn)
	}
}

func benchmarkHotPotatoCSnowflake(b *testing.B) {
	sf := snowflake.New(1)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sf.NextID()
	}
}

func benchmarkBwmarrinSnowflake(b *testing.B) {
	node, _ := bwmarrinsnowflake.NewNode(1)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		node.Generate()
	}
}

func benchmarkGoSnowflake(b *testing.B) {
	gosnowflake.SetMachineID(1)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		gosnowflake.ID()
	}
}
