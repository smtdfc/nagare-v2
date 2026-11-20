//go:build !dix

//go:generate dix wire
package main

import (
	"fmt"

	"github.com/smtdfc/nagare/core/.dix/generated"
)

func main() {
	entry, _ := generated.Root()
	err := entry.Start()
	if err != nil {
		fmt.Println("Error starting entry:", err)
		return
	}
}
