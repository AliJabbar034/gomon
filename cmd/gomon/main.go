package main

import (
	"log"
	"os"

	"github.com/alijabbar034/gomon/internal/watcher"
)

func main() {
    if len(os.Args) < 2 {
        log.Fatalf("Usage: gomon <command> [args]")
    }

    command := os.Args[1]
    args := os.Args[2:]

    // Start the watcher
    if err := watcher.Start(command, args); err != nil {
        log.Fatal(err)
    }
}
