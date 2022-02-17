package errors

import "io"

func MaybePanicOnClose(c io.Closer) {
	if c != nil {
		MaybePanic(c.Close())
	}
}

func MaybePanic(err error) {
	if err != nil {
		panic(err)
	}
}
