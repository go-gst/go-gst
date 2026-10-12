package gstsdp

import "iter"

// getIter is a utility function that creates an iterator from a length and an accessor function.
func getIter[T any](length uint32, accessor func(i uint32) T) iter.Seq2[uint32, T] {
	return func(yield func(uint32, T) bool) {
		for i := range length {
			if !yield(i, accessor(i)) {
				return
			}
		}
	}
}
