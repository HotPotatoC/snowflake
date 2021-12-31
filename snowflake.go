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
	// epoch defaults to March 3, 2012, 00:00:00 UTC. Which is
	// the release date of Go 1.0.
	// You can customize it by calling SetEpoch()
	//	// Example
	//	err := snowflake.SetEpoch(time.Date(2009, 1, 1, 0, 0, 0, 0, time.UTC))
	epoch = time.Date(2012, 3, 28, 0, 0, 0, 0, time.UTC)

	// ErrEpochIsZero is returned when the epoch is set to the zero time.
	ErrEpochIsZero = errors.New("epoch is zero")
	// ErrEpochFuture is returned when the epoch is in the future.
	ErrEpochFuture = errors.New("epoch is in the future")
)

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
}

// New returns a new snowflake.ID (max discriminator value: 1023)
func New(discriminator uint64) *ID {
	return &ID{discriminator: discriminator}
}

// NextID returns a new snowflake ID.
//
//	Format:
//	1011001001101101011001010111100000001011111111111000000000001
//	|--------------timestamp--------------|--disc---|----seq----|
func (id *ID) NextID() uint64 {
	id.mtx.Lock()
	defer id.mtx.Unlock()

	nowSinceEpoch := time.Since(epoch).Nanoseconds() / int64(time.Millisecond)

	// If this is the first call to NextID(), initialize the elapsedTime
	// and set sequence to zero.
	if id.elapsedTime < nowSinceEpoch {
		id.elapsedTime = nowSinceEpoch
		id.sequence = 0
	} else { // Otherwise, increment the sequence number.
		id.sequence = (id.sequence + 1) & maxSeqBits
		// if we've used up all the bits in the sequence number,
		// we need to increment the timestamp
		if id.sequence == 0 {
			id.elapsedTime++
		}
	}

	timestampSegment := uint64(nowSinceEpoch << (sequenceBits + discriminatorBits))
	discriminatorSegment := uint64(id.discriminator) << sequenceBits
	sequenceSegment := uint64(id.sequence)

	// if the discriminator is bigger than the max, we need to reset it
	if id.discriminator > maxDiscriminatorBits {
		discriminatorSegment = 0
	}

	return timestampSegment | discriminatorSegment | sequenceSegment
}

type SID struct {
	Timestamp     int64
	Sequence      uint64
	Discriminator uint64
}

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

	nowSinceEpoch := time.Since(epoch).Nanoseconds() / int64(time.Millisecond)

	// If this is the first call to NextID(), initialize the elapsedTime
	// and set sequence to zero.
	if id.elapsedTime < nowSinceEpoch {
		id.elapsedTime = nowSinceEpoch
		id.sequence = 0
	} else { // Otherwise, increment the sequence number.
		id.sequence = (id.sequence + 1) & maxSeqBits
		// if we've used up all the bits in the sequence number,
		// we need to increment the timestamp
		if id.sequence == 0 {
			id.elapsedTime++
		}
	}

	timestampSegment := uint64(nowSinceEpoch << (sequenceBits + discriminatorBits))
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

type SID2 struct {
	Timestamp      int64
	Sequence       uint64
	Discriminator1 uint64
	Discriminator2 uint64
}

func Parse2(sid uint64) SID2 {
	return SID2{
		Timestamp:      getTimestamp(sid),
		Sequence:       getSequence(sid),
		Discriminator1: getFirstDiscriminant(sid),
		Discriminator2: getSecondDiscriminant(sid),
	}
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
	shift := sequenceBits + discriminatorBits
	return int64(id>>shift) + epoch.UnixNano()/int64(time.Millisecond)
}

// getSequence returns the sequence number of a snowflake ID. (internal-use only)
func getSequence(id uint64) uint64 { return uint64(int(id) & maxSeqBits) }
