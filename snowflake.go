package snowflake

import (
	"errors"
	"sync"
	"time"
)

const (
	fieldBits        = 10
	sequenceBits     = 12
	maxFieldBits     = 0x3FF // 0x3FF shorthand for (1 << fieldBits) - 1 or 1023
	maxFieldHalfBits = 0x1F  // 0x1F shorthand for (1 << (fieldBits / 2)) - 1 or 31
	maxSeqBits       = 0xFFF // 0xFFF shorthand for (1 << sequenceBits) - 1 or 4095
)

var (
	epoch = time.Date(2012, 3, 28, 0, 0, 0, 0, time.UTC)

	// ErrEpochIsZero is returned when the epoch is set to the zero time.
	ErrEpochIsZero = errors.New("epoch is zero")
	// ErrEpochFuture is returned when the epoch is in the future.
	ErrEpochFuture = errors.New("epoch is in the future")
)

// Epoch returns the current configured epoch.
// epoch defaults to March 3, 2012, 00:00:00 UTC. Which is
// the release date of Go 1.0.
// You can customize it by calling SetEpoch()
//
//	// Example setting the epoch to 2010-01-01 00:00:00 UTC
//	err := snowflake.SetEpoch(time.Date(2010, 1, 1, 0, 0, 0, 0, time.UTC))
func Epoch() time.Time { return epoch }

// SetEpoch changes the epoch / starting time to a custom time.
// Can return 2 errors: ErrEpochIsZero and ErrEpochFuture.
// If the epoch is set to the zero time, ErrEpochIsZero is returned.
// If the epoch is in the future, ErrEpochFuture is returned.
func SetEpoch(e time.Time) error {
	e = e.UTC()

	if e.IsZero() {
		return ErrEpochIsZero
	}

	if e.After(time.Now().UTC()) {
		return ErrEpochFuture
	}

	epoch = e

	return nil
}

// ID is a custom type for a snowflake ID.
type ID struct {
	mtx         sync.Mutex
	field       uint64
	sequence    uint64
	elapsedTime int64
	lastID      uint64
}

// New returns a new snowflake.ID (max field value: 1023)
func New(field uint64) *ID {
	return &ID{field: field, lastID: 0}
}

// NextID returns a new snowflake ID.
//
//	Format:
//	1011001001101101011001010111100000001011111111111000000000001
//	|--------------timestamp--------------|--disc---|----seq----|
func (id *ID) NextID() uint64 {
	id.mtx.Lock()
	defer id.mtx.Unlock()

	nowSinceEpoch := msSinceEpoch()

	// reference: https://github.com/twitter-archive/snowflake/blob/snowflake-2010/src/main/scala/com/twitter/service/snowflake/IdWorker.scala#L81
	if nowSinceEpoch == id.elapsedTime { // same millisecond as last time
		id.sequence = (id.sequence + 1) & maxSeqBits // increment sequence number

		if id.sequence == 0 {
			// if we've used up all the bits in the sequence number,
			// we need to change the timestamp
			nowSinceEpoch = waitUntilNextMs(id.elapsedTime) // wait until next millisecond
		}
	} else {
		id.sequence = 0
	}

	id.elapsedTime = nowSinceEpoch

	timestampSegment := uint64(id.elapsedTime << (sequenceBits + fieldBits))
	fieldSegment := uint64(id.field) << sequenceBits
	sequenceSegment := uint64(id.sequence)

	// if the field is bigger than the max, we need to reset it
	if id.field > maxFieldBits {
		fieldSegment = 0
	}

	return timestampSegment | fieldSegment | sequenceSegment
}

// SID is the parsed representation of a snowflake ID.
type SID struct {
	// Timestamp is the timestamp of the snowflake ID.
	Timestamp int64
	// Sequence is the sequence number of the snowflake ID.
	Sequence uint64
	// Field is the field value of the snowflake ID.
	Field uint64
}

// Parse parses an existing snowflake ID
func Parse(sid uint64) SID {
	return SID{
		Timestamp: getTimestamp(sid),
		Sequence:  getSequence(sid),
		Field:     getDiscriminant(sid),
	}
}

// ID2 is a snowflake ID with 2 field fields.
type ID2 struct {
	mtx         sync.Mutex
	field1      uint64
	field2      uint64
	sequence    uint64
	elapsedTime int64
}

// New2 returns a new snowflake.ID2 (max field value: 31)
func New2(field1 uint64, field2 uint64) *ID2 {
	return &ID2{field1: field1, field2: field2}
}

// NextID returns a new snowflake ID with 2 field fields.
// The field fields are split into 5 bits each. (max field each: 31)
//
//	Format:
//	1011001001101101011001010111100000001011111111111000000000001
//	|--------------timestamp--------------|-d2-|-d1-|----seq----|
func (id *ID2) NextID() uint64 {
	id.mtx.Lock()
	defer id.mtx.Unlock()

	nowSinceEpoch := msSinceEpoch()

	// reference: https://github.com/twitter-archive/snowflake/blob/snowflake-2010/src/main/scala/com/twitter/service/snowflake/IdWorker.scala#L81
	if nowSinceEpoch == id.elapsedTime { // same millisecond as last time
		id.sequence = (id.sequence + 1) & maxSeqBits // increment sequence number

		if id.sequence == 0 {
			// if we've used up all the bits in the sequence number,
			// we need to change the timestamp
			nowSinceEpoch = waitUntilNextMs(id.elapsedTime) // wait until next millisecond
		}
	} else {
		id.sequence = 0
	}

	id.elapsedTime = nowSinceEpoch

	timestampSegment := uint64(id.elapsedTime << (sequenceBits + fieldBits))
	field1Segment := id.field1 << uint64(sequenceBits)
	field2Segment := id.field2 << uint64(sequenceBits+fieldBits/2)
	sequenceSegment := uint64(id.sequence)

	// if the field is bigger than the max, we need to reset it
	if id.field1 > uint64(maxFieldHalfBits) {
		field1Segment = 0
	}

	if id.field2 > uint64(maxFieldHalfBits) {
		field2Segment = 0
	}

	return timestampSegment | field2Segment | field1Segment | sequenceSegment
}

// SID2 is the parsed representation of a snowflake ID with 2 field fields.
type SID2 struct {
	// Timestamp is the timestamp of the snowflake ID.
	Timestamp int64
	// Sequence is the sequence number of the snowflake ID.
	Sequence uint64
	// Field1 is the first field value of the snowflake ID.
	Field1 uint64
	// Field2 is the second field value of the snowflake ID.
	Field2 uint64
}

// Parse2 parses an existing snowflake ID with 2 field fields.
func Parse2(sid uint64) SID2 {
	return SID2{
		Timestamp: getTimestamp(sid),
		Sequence:  getSequence(sid),
		Field1:    getFirstDiscriminant(sid),
		Field2:    getSecondDiscriminant(sid),
	}
}

// waitUntilNextMs waits until the next millisecond to return. (internal-use only)
func waitUntilNextMs(last int64) int64 {
	ms := msSinceEpoch()
	for ms <= last {
		ms = msSinceEpoch()
	}
	return ms
}

// msSinceEpoch returns the number of milliseconds since the epoch. (internal-use only)
func msSinceEpoch() int64 {
	return time.Since(epoch).Nanoseconds() / 1e6
}

// getDiscriminant returns the discriminant value of a snowflake ID. (internal-use only)
func getDiscriminant(id uint64) uint64 {
	return (id >> sequenceBits) & maxFieldBits
}

// getFirstDiscriminant returns the first discriminant value of a snowflake ID. (internal-use only)
func getFirstDiscriminant(id uint64) uint64 {
	return (id >> sequenceBits) & maxFieldHalfBits
}

// getSecondDiscriminant returns the second discriminant value of a snowflake ID. (internal-use only)
func getSecondDiscriminant(id uint64) uint64 {
	return (id >> (sequenceBits + fieldBits/2)) & maxFieldHalfBits
}

// getTimestamp returns the timestamp of a snowflake ID. (internal-use only)
func getTimestamp(id uint64) int64 {
	return int64(id>>(sequenceBits+fieldBits)) + epoch.UnixNano()/1e6
}

// getSequence returns the sequence number of a snowflake ID. (internal-use only)
func getSequence(id uint64) uint64 { return uint64(int(id) & maxSeqBits) }
