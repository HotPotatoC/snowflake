package snowflake_test

import (
	"testing"
	"time"

	"github.com/HotPotatoC/snowflake"
)

func TestSequence(t *testing.T) {
	sf := snowflake.New(1)

	id := sf.NextID()
	if snowflake.Parse(id).Sequence != 0 {
		t.Errorf("expected sequence 0 got %d", snowflake.Parse(id).Sequence)
	}

	id = sf.NextID()
	if snowflake.Parse(id).Sequence != 1 {
		t.Errorf("expected sequence 1 got %d", snowflake.Parse(id).Sequence)
	}

	id = sf.NextID()
	if snowflake.Parse(id).Sequence != 2 {
		t.Errorf("expected sequence 2 got %d", snowflake.Parse(id).Sequence)
	}

	time.Sleep(time.Millisecond * 100)

	id = sf.NextID()
	if snowflake.Parse(id).Sequence != 0 {
		t.Errorf("expected sequence to return to 0 got %d", snowflake.Parse(id).Sequence)
	}
}

func TestSequence2(t *testing.T) {
	sf := snowflake.New2(1, 1)

	id := sf.NextID()
	if snowflake.Parse(id).Sequence != 0 {
		t.Errorf("expected sequence 0 got %d", snowflake.Parse(id).Sequence)
	}

	id = sf.NextID()
	if snowflake.Parse(id).Sequence != 1 {
		t.Errorf("expected sequence 1 got %d", snowflake.Parse(id).Sequence)
	}

	id = sf.NextID()
	if snowflake.Parse(id).Sequence != 2 {
		t.Errorf("expected sequence 2 got %d", snowflake.Parse(id).Sequence)
	}

	time.Sleep(time.Millisecond * 100)

	id = sf.NextID()
	if snowflake.Parse(id).Sequence != 0 {
		t.Errorf("expected sequence to return to 0 got %d", snowflake.Parse(id).Sequence)
	}
}

func TestDiscriminator(t *testing.T) {
	tc := []struct {
		name     string
		sf       *snowflake.ID
		expected uint64
	}{
		{"1", snowflake.New(1), 1},
		{"100", snowflake.New(100), 100},
		{"1000", snowflake.New(1000), 1000},
		{"1023", snowflake.New(1023), 1023},
		{"1024 Overflow", snowflake.New(1024), 0},
	}

	for _, tt := range tc {
		t.Run(tt.name, func(t *testing.T) {
			if snowflake.Parse(tt.sf.NextID()).Discriminator != tt.expected {
				t.Errorf("expected discriminator %d got %d",
					tt.expected, snowflake.Parse(tt.sf.NextID()).Discriminator)
			}
		})
	}
}

func Test2Discriminators(t *testing.T) {
	tc := []struct {
		name     string
		sf       *snowflake.ID2
		expected [2]uint64
	}{
		{"1", snowflake.New2(1, 1), [2]uint64{1, 1}},
		{"10", snowflake.New2(10, 10), [2]uint64{10, 10}},
		{"31", snowflake.New2(31, 31), [2]uint64{31, 31}},
		{"32 overflow", snowflake.New2(32, 32), [2]uint64{0, 0}},
	}

	for _, tt := range tc {
		t.Run(tt.name, func(t *testing.T) {
			firstDiscriminator, secondDiscriminator := tt.expected[0], tt.expected[1]
			parsed := snowflake.Parse2(tt.sf.NextID())
			if parsed.Discriminator1 != firstDiscriminator {
				t.Errorf("expected discriminator %d got %d",
					tt.expected, parsed.Discriminator1)
			}
			if parsed.Discriminator2 != secondDiscriminator {
				t.Errorf("expected 2nd discriminator %d got %d", secondDiscriminator, parsed.Discriminator2)
			}
		})
	}
}
