package common

func Must[T any](v T, err error) T {
	if err != nil {
		panic("got error:" + err.Error())
	}

	return v
}
