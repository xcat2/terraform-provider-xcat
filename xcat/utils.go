package xcat

import (
//"fmt"
//"strings"
)

func Contains(a []string, x string) bool {
	for _, n := range a {
		if x == n {
			return true
		}
	}
	return false

}
