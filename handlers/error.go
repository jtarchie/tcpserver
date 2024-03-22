package handlers

import (
	"bufio"
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"

	"github.com/jtarchie/tcpserver"
)

type Error struct{}

var _ tcpserver.Handler = &Error{}

var ErrOnConnection = errors.New("this always occurs")

func (*Error) OnConnection(_ context.Context, input io.ReadWriter) error {
	reader := bufio.NewReader(input)

	for {
		line, _, err := reader.ReadLine()
		if err != nil {
			return fmt.Errorf("unexpected error: %w", err)
		}

		if bytes.Equal(line, []byte(`error`)) {
			return fmt.Errorf("something happened: %w", ErrOnConnection)
		}

		_, _ = input.Write(line)
		_, _ = input.Write([]byte("\r\n"))
	}
}
