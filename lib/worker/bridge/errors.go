package bridge

import "errors"

var errIsFull = errors.New("bridge is full")

func IsFullErr(err error) bool {
	return errors.Is(err, errIsFull)
}
