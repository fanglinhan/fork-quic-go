package congestion

import (
	"github.com/quic-go/quic-go/internal/protocol"
	"github.com/quic-go/quic-go/internal/utils/ringbuffer"
)

// PacketNumberIndexedQueue is a queue of mostly continuous numbered entries
// which supports the following operations:
// - adding elements to the end of the queue, or at some point past the end
// - removing elements in any order
// - retrieving elements
// If all elements are inserted in order, all of the operations above are
// amortized O(1) time.
//
// Internally, the data structure is a deque where each element is marked as
// present or not.  The deque starts at the lowest present index.  Whenever an
// element is removed, it's marked as not present, and the front of the deque is
// cleared of elements that are not present.
//
// The tail of the queue is not cleared due to the assumption of entries being
// inserted in order, though removing all elements of the queue will return it
// to its initial state.
//
// Note that this data structure is inherently hazardous, since an addition of
// just two entries will cause it to consume all of the memory available.
// Because of that, it is not a general-purpose container and should not be used
// as one.

type entryWrapper[T any] struct {
	present bool
	entry   T
}

type PacketNumberIndexedQueue[T any] struct {
	entries                ringbuffer.RingBuffer[entryWrapper[T]]
	numberOfPresentEntries int
	firstPacket            protocol.PacketNumber
}

func NewPacketNumberIndexedQueue[T any](size int) *PacketNumberIndexedQueue[T] {
	q := &PacketNumberIndexedQueue[T]{
		firstPacket: protocol.InvalidPacketNumber,
	}

	q.entries.Init(size)

	return q
}

// Emplace inserts data associated |packet_number| into (or past) the end of the
// queue, filling up the missing intermediate entries as necessary.  Returns
// true if the element has been inserted successfully, false if it was already
// in the queue or inserted out of order.
func (p *PacketNumberIndexedQueue[T]) Emplace(packetNumber protocol.PacketNumber, entry *T) bool {
	if packetNumber == protocol.InvalidPacketNumber || entry == nil {
		return false
	}

	if p.IsEmpty() {
		p.entries.PushBack(entryWrapper[T]{
			present: true,
			entry:   *entry,
		})
		p.numberOfPresentEntries = 1
		p.firstPacket = packetNumber
		return true
	}

	// Do not allow insertion out-of-order.
	if packetNumber <= p.LastPacket() {
		return false
	}

	// Handle potentially missing elements.
	offset := int(packetNumber - p.FirstPacket())
	if gap := offset - p.entries.Len(); gap > 0 {
		for i := 0; i < gap; i++ {
			p.entries.PushBack(entryWrapper[T]{})
		}
	}

	p.entries.PushBack(entryWrapper[T]{
		present: true,
		entry:   *entry,
	})
	p.numberOfPresentEntries++
	return true
}

// GetEntry Retrieve the entry associated with the packet number.  Returns the pointer
// to the entry in case of success, or nullptr if the entry does not exist.
func (p *PacketNumberIndexedQueue[T]) GetEntry(packetNumber protocol.PacketNumber) (entry *T) {
	return nil
}

// Remove, Same as above, but if an entry is present in the queue, also call f(entry)
// before removing it.
func (p *PacketNumberIndexedQueue[T]) Remove(packetNumber protocol.PacketNumber) bool {
	return false
}

// RemoveUpTo, but not including |packet_number|.
// Unused slots in the front are also removed, which means when the function
// returns, |first_packet()| can be larger than |packet_number|.
func (p *PacketNumberIndexedQueue[T]) RemoveUpTo(packetNumber protocol.PacketNumber) {
	return
}

// IsEmpty return if queue is empty.
func (p *PacketNumberIndexedQueue[T]) IsEmpty() bool {
	return p.numberOfPresentEntries == 0
}

// NumberOfPresentEntries returns the number of entries in the queue.
func (p *PacketNumberIndexedQueue[T]) NumberOfPresentEntries(packetNumber protocol.PacketNumber) int {
	return p.numberOfPresentEntries
}

// EntrySlotsUsed returns the number of entries allocated in the underlying deque.  This is
// proportional to the memory usage of the queue.
func (p *PacketNumberIndexedQueue[T]) EntrySlotsUsed(packetNumber protocol.PacketNumber) int {
	return p.entries.Len()
}

// LastPacket returns packet number of the first entry in the queue.
func (p *PacketNumberIndexedQueue[T]) FirstPacket() (packetNumber protocol.PacketNumber) {
	return p.firstPacket
}

// LastPacket returns packet number of the last entry ever inserted in the queue.  Note that the
// entry in question may have already been removed.  Zero if the queue is
// empty.
func (p *PacketNumberIndexedQueue[T]) LastPacket() (packetNumber protocol.PacketNumber) {
	if p.IsEmpty() {
		return protocol.InvalidPacketNumber
	}

	return p.firstPacket + protocol.PacketNumber(p.entries.Len()-1)
}
