//go:build !dix

//go:generate dix wire
package main

import (
	"fmt"

	"github.com/smtdfc/nagare/core/.dix/generated"
)

func main() {
	entry, err := generated.Root()
	if err != nil {
		fmt.Printf("Error %s", err)
		return
	}

	if entry == nil {
		return
	}

	err = entry.Start()
	if err != nil {
		entry.Logger.Error("Error starting entry", "error", err)
		return
	}
}
