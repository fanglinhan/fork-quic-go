package congestion

import (
	"math"
	"time"

	"github.com/quic-go/quic-go/internal/protocol"
)

const (
	infBandwidth = Bandwidth(math.MaxUint64)
	infRTT       = time.Duration(math.MaxInt64)
)

// Bandwidth of a connection
type Bandwidth uint64

const (
	// BitsPerSecond is 1 bit per second
	BitsPerSecond Bandwidth = 1
	// BytesPerSecond is 1 byte per second
	BytesPerSecond = 8 * BitsPerSecond
)

// BandwidthFromDelta calculates the bandwidth from a number of bytes and a time delta
func BandwidthFromDelta(bytes protocol.ByteCount, delta time.Duration) Bandwidth {
	return Bandwidth(bytes) * Bandwidth(time.Second) / Bandwidth(delta) * BytesPerSecond
}

// BytesFromBandwidthAndTimeDelta calculates the bytes
// from a bandwidth(bits per second) and a time delta
func BytesFromBandwidthAndTimeDelta(bandwidth Bandwidth, delta time.Duration) protocol.ByteCount {
	return (protocol.ByteCount(bandwidth) * protocol.ByteCount(delta)) /
		(protocol.ByteCount(time.Second) * 8)
}
