package watcher

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"sync"

	"github.com/fsnotify/fsnotify"
)

var (
    cmdMutex sync.Mutex
    cmd      *exec.Cmd
)

// Start initializes the file watcher and handles the restart logic
func Start(command string, args []string) error {
    watcher, err := fsnotify.NewWatcher()
    if err != nil {
        return err
    }
    defer watcher.Close()

    go watchForChanges(watcher, command, args)

    err = filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
        if err != nil {
            return err
        }
        if info.IsDir() {
            return watcher.Add(path)
        }
        return nil
    })
    if err != nil {
        return err
    }

    // Initial run
    restart(command, args)

    done := make(chan bool)
    <-done

    return nil
}

func watchForChanges(watcher *fsnotify.Watcher, command string, args []string) {
    for {
        select {
        case event, ok := <-watcher.Events:
            if !ok {
                return
            }
            fmt.Println("File changed:", event)
            if event.Op&fsnotify.Write == fsnotify.Write {
                restart(command, args)
            }
        case err, ok := <-watcher.Errors:
            if !ok {
                return
            }
            log.Println("error:", err)
        }
    }
}

func restart(command string, args []string) {
    cmdMutex.Lock()
    defer cmdMutex.Unlock()

    if cmd != nil {
        cmd.Process.Kill()
        cmd.Wait()
    }

    cmd = exec.Command(command, args...)
    cmd.Stdout = os.Stdout
    cmd.Stderr = os.Stderr
    go cmd.Run()
}
