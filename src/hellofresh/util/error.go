// Package util utilities including auth, reponse, error handling
package util

import "log"

// PanicOnError log on panic when error
func PanicOnError(err error) {
	log.Fatal(err)
	panic(err)
}
