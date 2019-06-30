package main

import (
	"fmt"
)

func UnsupportedError(errMsg string) error {
	return fmt.Errorf("UnsupportedError:%s", errMsg)
}
