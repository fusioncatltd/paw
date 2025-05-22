package utils

import (
	"bytes"
	"context"
	"github.com/urfave/cli/v3"
	"io"
	"log"
	"os"
)

func CaptureOutputInTests(f func(context.Context, *cli.Command) error, ctx context.Context, cmd *cli.Command) (string, error) {
	// 1) keep a reference to the real stdout
	oldStdout := os.Stdout

	// 2) create a pipe
	r, w, err := os.Pipe()
	if err != nil {
		panic("could not create pipe: " + err.Error())
	}

	// 3) redirect stdout to the pipe writer
	os.Stdout = w

	// run the function
	err = f(ctx, cmd)

	// 4) close writer, restore stdout
	w.Close()
	os.Stdout = oldStdout

	// 5) read the captured output
	var buf bytes.Buffer
	if _, err := io.Copy(&buf, r); err != nil {
		panic("could not read captured output: " + err.Error())
	}
	r.Close()

	return buf.String(), err
}

func Cleanup() {
	tempDir, err := os.MkdirTemp("", "paw-test-*")
	if err != nil {
		log.Fatal("Failed to create temp dir: %v", err)
	}

	// Remove all contents of the temp directory
	if err := os.RemoveAll(tempDir); err != nil {
		log.Fatal("Failed to clean temp dir: %v", err)
	}
	// Recreate the empty temp directory
	if err := os.Mkdir(tempDir, 0755); err != nil {
		log.Fatal("Failed to recreate temp dir: %v", err)
	}
	// Change back to the temp directory
	if err := os.Chdir(tempDir); err != nil {
		log.Fatal("Failed to change to temp dir: %v", err)
	}
}
