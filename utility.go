// utility.go
package main

import "log"

// handleErrorMessage is a helper function to log fatal errors.
func handleErrorMessage(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}
