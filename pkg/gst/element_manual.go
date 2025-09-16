package gst

type ElementExtManual interface {
	BlockSetState(state State, timeout ClockTime) StateChangeReturn
}

// BlockSetState is a convenience wrapper around calling SetState and State to wait for async state changes. See State for more info.
func (el *ElementInstance) BlockSetState(state State, timeout ClockTime) StateChangeReturn {
	ret := el.SetState(state)

	if ret == StateChangeAsync {
		_, _, ret = el.GetState(timeout)
	}

	return ret
}

func LinkMany(elements ...Element) bool {
	if len(elements) < 2 {
		return false
	}

	for i := 0; i < len(elements)-1; i++ {
		if !elements[i].Link(elements[i+1]) {
			return false
		}
	}

	return true
}
