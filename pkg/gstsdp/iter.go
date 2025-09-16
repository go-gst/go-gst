package gstsdp

import "iter"

// getIter is a utility function that creates an iterator from a length and an accessor function.
func getIter[T any](length uint, accessor func(i uint) T) iter.Seq2[uint, T] {
	return func(yield func(uint, T) bool) {
		for i := uint(0); i < length; i++ {
			if !yield(i, accessor(i)) {
				return
			}
		}
	}
}
