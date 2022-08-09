package snowflake_test

import (
	"sync"
	"testing"
	"time"

	"github.com/HotPotatoC/snowflake"
)

func TestNextID(t *testing.T) {
	n := 100000
	sf := snowflake.New(1)
	ids := make(map[uint64]bool)
	for i := 0; i < n; i++ {
		id := sf.NextID()
		if _, exists := ids[id]; exists {
			t.Errorf("expected to be unique, but got a repeated ID (%d)", id)
			break
		}

		ids[id] = true
	}
}

func TestNextID_Concurrent(t *testing.T) {
	n := 100000
	ch := make(chan uint64, n)
	sf := snowflake.New(1)

	var wg sync.WaitGroup
	for i := 0; i < n; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			id := sf.NextID()
			ch <- id
		}()
	}
	wg.Wait()
	close(ch)

	ids := make(map[uint64]bool)
	for id := range ch {
		if _, ok := ids[id]; ok {
			t.Error("expected to be unique, but got a repeated ID")
			break
		}
		ids[id] = true
	}
	if len(ids) != n {
		t.Errorf("expected map length %d got %d", n, len(ids))
	}
}

func TestNextID2(t *testing.T) {
	n := 100000
	sf := snowflake.New2(1, 1)
	ids := make(map[uint64]bool)
	for i := 0; i < n; i++ {
		id := sf.NextID()
		if _, exists := ids[id]; exists {
			t.Error("expected to be unique, but got a repeated ID")
			break
		}

		ids[id] = true
	}
}

func TestNextID2_Concurrent(t *testing.T) {
	n := 100000
	ch := make(chan uint64, n)
	sf := snowflake.New2(1, 1)

	var wg sync.WaitGroup
	for i := 0; i < n; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			id := sf.NextID()
			ch <- id
		}()
	}
	wg.Wait()
	close(ch)

	ids := make(map[uint64]bool)
	for id := range ch {
		if _, ok := ids[id]; ok {
			t.Error("expected to be unique, but got a repeated ID")
			break
		}
		ids[id] = true
	}
	if len(ids) != n {
		t.Errorf("expected map length %d got %d", n, len(ids))
	}
}

func TestID_IsIncreasing(t *testing.T) {
	sf := snowflake.New(1)
	n := 100000
	ids := make([]uint64, n)
	for i := 0; i < n; i++ {
		ids = append(ids, sf.NextID())
	}

	for i := 0; i < n; i++ {
		if i > 0 && ids[i] < ids[i-1] {
			t.Errorf("expected to be increasing, but got %d at %d", ids[i], i)
			break
		}
	}
}

func TestID2_IsIncreasing(t *testing.T) {
	sf := snowflake.New2(1, 2)
	n := 100000
	ids := make([]uint64, n)
	for i := 0; i < n; i++ {
		ids = append(ids, sf.NextID())
	}

	for i := 0; i < n; i++ {
		if i > 0 && ids[i] < ids[i-1] {
			t.Errorf("expected to be increasing, but got %d at %d", ids[i], i)
			break
		}
	}
}

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

func TestField(t *testing.T) {
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
			if snowflake.Parse(tt.sf.NextID()).Field != tt.expected {
				t.Errorf("expected field %d got %d",
					tt.expected, snowflake.Parse(tt.sf.NextID()).Field)
			}
		})
	}
}

func Test2Fields(t *testing.T) {
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
			firstField, secondField := tt.expected[0], tt.expected[1]
			parsed := snowflake.Parse2(tt.sf.NextID())
			if parsed.Field1 != firstField {
				t.Errorf("expected field %d got %d",
					tt.expected, parsed.Field1)
			}
			if parsed.Field2 != secondField {
				t.Errorf("expected 2nd field %d got %d", secondField, parsed.Field2)
			}
		})
	}
}

func TestParse(t *testing.T) {
	// timestamp: 1640942460724
	// Field: 1
	// Sequence: 0
	id := uint64(1292053924173320192)

	sid := snowflake.Parse(id)
	if sid.Sequence != 0 {
		t.Errorf("expected sequence %d got %d", 0, sid.Sequence)
	}

	if sid.Field != 1 {
		t.Errorf("expected field %d got %d", 1, sid.Field)
	}

	if sid.Timestamp != 1640942460724 {
		t.Errorf("expected timestamp %d got %d", 1640942460724, sid.Timestamp)
	}
}

func TestParse2Fields(t *testing.T) {
	// timestamp: 1640945127245
	// Field1: 1
	// Field2: 24
	// Sequence: 0
	id := uint64(1292065108376162304)

	sid := snowflake.Parse2(id)
	if sid.Sequence != 0 {
		t.Errorf("expected sequence %d got %d", 0, sid.Sequence)
	}

	if sid.Field1 != 1 {
		t.Errorf("expected field %d got %d", 1, sid.Field1)
	}

	if sid.Field2 != 24 {
		t.Errorf("expected field %d got %d", 24, sid.Field1)
	}

	if sid.Timestamp != 1640945127245 {
		t.Errorf("expected timestamp %d got %d", 1640945127245, sid.Timestamp)
	}
}

func TestEpoch(t *testing.T) {
	epoch := time.Date(2012, 3, 28, 0, 0, 0, 0, time.UTC)

	if snowflake.Epoch() != epoch {
		t.Errorf("expected epoch %s got %s", epoch, snowflake.Epoch())
	}
}

func TestSetEpoch(t *testing.T) {
	defaultEpoch := time.Date(2012, 3, 28, 0, 0, 0, 0, time.UTC)
	tc := []struct {
		name     string
		epoch    time.Time
		expected time.Time
		err      error
	}{
		{"Should return ErrEpochIsZero", time.Time{}, defaultEpoch, snowflake.ErrEpochIsZero},
		{"Should return ErrEpochFuture", time.Date(3000, 1, 1, 0, 0, 0, 0, time.UTC), defaultEpoch, snowflake.ErrEpochFuture},
		{"2010-1-1 00:00:00", time.Date(2010, 1, 1, 0, 0, 0, 0, time.UTC), time.Date(2010, 1, 1, 0, 0, 0, 0, time.UTC), nil},
	}

	for _, tt := range tc {
		t.Run(tt.name, func(t *testing.T) {
			err := snowflake.SetEpoch(tt.epoch)

			if snowflake.Epoch() != tt.expected {
				t.Errorf("expected epoch %s got %s", tt.expected, snowflake.Epoch())
			}

			if err != tt.err {
				t.Errorf("expected error %v got %v", tt.err, err)
			}
		})
	}
}
