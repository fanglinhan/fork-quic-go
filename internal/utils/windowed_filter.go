package utils

import (
	"golang.org/x/exp/constraints"
)

// Implements Kathleen Nichols' algorithm for tracking the minimum (or maximum)
// estimate of a stream of samples over some fixed time interval. (E.g.,
// the minimum RTT over the past five minutes.) The algorithm keeps track of
// the best, second best, and third best min (or max) estimates, maintaining an
// invariant that the measurement time of the n'th best >= n-1'th best.

// The algorithm works as follows. On a reset, all three estimates are set to
// the same sample. The second best estimate is then recorded in the second
// quarter of the window, and a third best estimate is recorded in the second
// half of the window, bounding the worst case error when the true min is
// monotonically increasing (or true max is monotonically decreasing) over the
// window.
//
// A new best sample replaces all three estimates, since the new best is lower
// (or higher) than everything else in the window and it is the most recent.
// The window thus effectively gets reset on every new min. The same property
// holds true for second best and third best estimates. Specifically, when a
// sample arrives that is better than the second best but not better than the
// best, it replaces the second and third best estimates but not the best
// estimate. Similarly, a sample that is better than the third best estimate
// but not the other estimates replaces only the third best estimate.
//
// Finally, when the best expires, it is replaced by the second best, which in
// turn is replaced by the third best. The newest sample replaces the third
// best.

type WindowedFilterValue interface {
	comparable
}

type WindowedFilterTime interface {
	constraints.Integer | constraints.Float
}

type WindowedFilter[V WindowedFilterValue, T WindowedFilterTime] struct {
	// Time length of window.
	windowLength T
	estimates    []Sample[V, T]
	comparator   func(V, V) bool
}

type Sample[V WindowedFilterValue, T WindowedFilterTime] struct {
	sample V
	time   T
}

// Compares two values and returns true if the first is greater than or equal
// to the second.
func MaxFilter[O constraints.Ordered](a, b O) bool {
	return a >= b
}

// Compares two values and returns true if the first is less than or equal
// to the second.
func MinFilter[O constraints.Ordered](a, b O) bool {
	return a <= b
}

func NewWindowedFilter[V WindowedFilterValue, T WindowedFilterTime](windowLength T, comparator func(V, V) bool) *WindowedFilter[V, T] {
	return &WindowedFilter[V, T]{
		windowLength: windowLength,
		estimates:    make([]Sample[V, T], 3, 3),
		comparator:   comparator,
	}
}

// Changes the window length.  Does not update any current samples.
func (f *WindowedFilter[V, T]) SetWindowLength(windowLength T) {
	f.windowLength = windowLength
}

func (f *WindowedFilter[V, T]) GetBest() V {
	return f.estimates[0].sample
}

func (f *WindowedFilter[V, T]) GetSecondBest() V {
	return f.estimates[1].sample
}

func (f *WindowedFilter[V, T]) GetThirdBest() V {
	return f.estimates[2].sample
}

// Updates best estimates with |sample|, and expires and updates best
// estimates as necessary.
func (f *WindowedFilter[V, T]) Update(sample V, time T) {
	// Reset all estimates if they have not yet been initialized, if new sample
	// is a new best, or if the newest recorded estimate is too old.
	if f.estimates[0].time == 0 ||
		f.comparator(sample, f.estimates[0].sample) ||
		time-f.estimates[2].time > f.windowLength {
		f.Reset(sample, time)
		return
	}

	if f.comparator(sample, f.estimates[1].sample) {
		f.estimates[1].sample = sample
		f.estimates[1].time = time
		f.estimates[2].sample = sample
		f.estimates[2].time = time
	} else if f.comparator(sample, f.estimates[2].sample) {
		f.estimates[2].sample = sample
		f.estimates[2].time = time
	}

	// Expire and update estimates as necessary.
	if time-f.estimates[0].time > f.windowLength {
		// The best estimate hasn't been updated for an entire window, so promote
		// second and third best estimates.
		f.estimates[0].sample = f.estimates[1].sample
		f.estimates[0].time = f.estimates[1].time
		f.estimates[1].sample = f.estimates[2].sample
		f.estimates[1].time = f.estimates[2].time
		f.estimates[2].sample = sample
		f.estimates[2].time = time
		// Need to iterate one more time. Check if the new best estimate is
		// outside the window as well, since it may also have been recorded a
		// long time ago. Don't need to iterate once more since we cover that
		// case at the beginning of the method.
		if time-f.estimates[0].time > f.windowLength {
			f.estimates[0].sample = f.estimates[1].sample
			f.estimates[0].time = f.estimates[1].time
			f.estimates[1].sample = f.estimates[2].sample
			f.estimates[1].time = f.estimates[2].time
		}
		return
	}
	if f.estimates[1].sample == f.estimates[0].sample && time-f.estimates[1].time > f.windowLength/4 {
		// A quarter of the window has passed without a better sample, so the
		// second-best estimate is taken from the second quarter of the window.
		f.estimates[1].sample = sample
		f.estimates[1].time = time
		f.estimates[2].sample = sample
		f.estimates[2].time = time
		return
	}

	if f.estimates[2].sample == f.estimates[1].sample && time-f.estimates[2].time > f.windowLength/2 {
		// We've passed a half of the window without a better estimate, so take
		// a third-best estimate from the second half of the window.
		f.estimates[2].sample = sample
		f.estimates[2].time = time
	}
}

// Resets all estimates to new sample.
func (f *WindowedFilter[V, T]) Reset(newSample V, newTime T) {
	f.estimates[0].sample = newSample
	f.estimates[0].time = newTime
	f.estimates[1].sample = newSample
	f.estimates[1].time = newTime
	f.estimates[2].sample = newSample
	f.estimates[2].time = newTime
}

func (f *WindowedFilter[V, T]) Clear() {
	f.estimates = make([]Sample[V, T], 3, 3)
}
