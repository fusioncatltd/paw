package utils

import (
	"bytes"
	"context"
	"github.com/urfave/cli/v3"
	"io"
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
