package snowflake

import (
	"errors"
	"sync"
	"time"
)

const (
	discriminatorBits        = 10
	sequenceBits             = 12
	maxDiscriminatorBits     = 0x3FF // 0x3FF shorthand for (1 << discriminatorBits) - 1 or 1023
	maxDiscriminatorHalfBits = 0x1F  // 0x1F shorthand for (1 << (discriminatorBits / 2)) - 1 or 31
	maxSeqBits               = 0xFFF // 0xFFF shorthand for (1 << sequenceBits) - 1 or 4095
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
	mtx           sync.Mutex
	discriminator uint64
	sequence      uint64
	elapsedTime   int64
	lastID        uint64
}

// New returns a new snowflake.ID (max discriminator value: 1023)
func New(discriminator uint64) *ID {
	return &ID{discriminator: discriminator, lastID: 0}
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

	timestampSegment := uint64(id.elapsedTime << (sequenceBits + discriminatorBits))
	discriminatorSegment := uint64(id.discriminator) << sequenceBits
	sequenceSegment := uint64(id.sequence)

	// if the discriminator is bigger than the max, we need to reset it
	if id.discriminator > maxDiscriminatorBits {
		discriminatorSegment = 0
	}

	return timestampSegment | discriminatorSegment | sequenceSegment
}

// SID is the parsed representation of a snowflake ID.
type SID struct {
	// Timestamp is the timestamp of the snowflake ID.
	Timestamp int64
	// Sequence is the sequence number of the snowflake ID.
	Sequence uint64
	// Discriminator is the discriminator value of the snowflake ID.
	Discriminator uint64
}

// Parse parses an existing snowflake ID
func Parse(sid uint64) SID {
	return SID{
		Timestamp:     getTimestamp(sid),
		Sequence:      getSequence(sid),
		Discriminator: getDiscriminant(sid),
	}
}

// ID2 is a snowflake ID with 2 discriminator fields.
type ID2 struct {
	mtx            sync.Mutex
	discriminator1 uint64
	discriminator2 uint64
	sequence       uint64
	elapsedTime    int64
}

// New2 returns a new snowflake.ID2 (max discriminator value: 31)
func New2(discriminator1 uint64, discriminator2 uint64) *ID2 {
	return &ID2{discriminator1: discriminator1, discriminator2: discriminator2}
}

// NextID returns a new snowflake ID with 2 discriminator fields.
// The discriminator fields are split into 5 bits each. (max discriminator each: 31)
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

	timestampSegment := uint64(id.elapsedTime << (sequenceBits + discriminatorBits))
	discriminator1Segment := id.discriminator1 << uint64(sequenceBits)
	discriminator2Segment := id.discriminator2 << uint64(sequenceBits+discriminatorBits/2)
	sequenceSegment := uint64(id.sequence)

	// if the discriminator is bigger than the max, we need to reset it
	if id.discriminator1 > uint64(maxDiscriminatorHalfBits) {
		discriminator1Segment = 0
	}

	if id.discriminator2 > uint64(maxDiscriminatorHalfBits) {
		discriminator2Segment = 0
	}

	return timestampSegment | discriminator2Segment | discriminator1Segment | sequenceSegment
}

// SID2 is the parsed representation of a snowflake ID with 2 discriminator fields.
type SID2 struct {
	// Timestamp is the timestamp of the snowflake ID.
	Timestamp int64
	// Sequence is the sequence number of the snowflake ID.
	Sequence uint64
	// Discriminator1 is the first discriminator value of the snowflake ID.
	Discriminator1 uint64
	// Discriminator2 is the second discriminator value of the snowflake ID.
	Discriminator2 uint64
}

// Parse2 parses an existing snowflake ID with 2 discriminator fields.
func Parse2(sid uint64) SID2 {
	return SID2{
		Timestamp:      getTimestamp(sid),
		Sequence:       getSequence(sid),
		Discriminator1: getFirstDiscriminant(sid),
		Discriminator2: getSecondDiscriminant(sid),
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
	return (id >> sequenceBits) & maxDiscriminatorBits
}

// getFirstDiscriminant returns the first discriminant value of a snowflake ID. (internal-use only)
func getFirstDiscriminant(id uint64) uint64 {
	return (id >> sequenceBits) & maxDiscriminatorHalfBits
}

// getSecondDiscriminant returns the second discriminant value of a snowflake ID. (internal-use only)
func getSecondDiscriminant(id uint64) uint64 {
	return (id >> (sequenceBits + discriminatorBits/2)) & maxDiscriminatorHalfBits
}

// getTimestamp returns the timestamp of a snowflake ID. (internal-use only)
func getTimestamp(id uint64) int64 {
	return int64(id>>(sequenceBits+discriminatorBits)) + epoch.UnixNano()/1e6
}

// getSequence returns the sequence number of a snowflake ID. (internal-use only)
func getSequence(id uint64) uint64 { return uint64(int(id) & maxSeqBits) }
