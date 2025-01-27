package main

import (
	"os"
	"testing"
	"time"
)

func Test_Main(t *testing.T) {
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()

	os.Args = []string{oldArgs[0], "-a=localhost:0"}

	done := make(chan struct{})

	go func() {
		main()
		close(done)
	}()

	time.Sleep(100 * time.Millisecond)

	p, err := os.FindProcess(os.Getpid())
	if err == nil && p != nil {
		_ = p.Signal(os.Interrupt)
	}

	select {
	case <-done:
		// main() exited successfully
	case <-time.After(1 * time.Second):
		t.Fatal("timeout: main() did not exit after signal")
	}
}
