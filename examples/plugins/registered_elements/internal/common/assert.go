package common

import "fmt"

var FinalizersCalled int = 0

func AssertFinalizersCalled(x int) {
	if FinalizersCalled != x {
		panic(fmt.Sprintf("finalizers did not run correctly, memory leak, wanted: %d, got: %d", x, FinalizersCalled))
	}
}
