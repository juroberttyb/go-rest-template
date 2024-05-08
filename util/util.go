package util

import (
	"errors"
)

func init() {
}
func ChErrHandler(errCh <-chan error, spawned int) error {
	for i := 0; i < spawned; i++ {
		select {
		case errVal, ok := <-errCh:
			if ok {
				if errVal != nil {
					return errVal
				}
			} else {
				return errors.New("error channel closed")
			}
		default:
		}
	}
	return nil
}
