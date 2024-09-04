package jsonl

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Writer struct {
	w      io.Writer
	prefix string
}

type WriterOption func(w *Writer)

func WithPrefix(prefix string) WriterOption {
	return func(w *Writer) {
		w.prefix = prefix
	}
}

func NewWriter(w io.Writer, opts ...WriterOption) Writer {
	wr := Writer{
		w: w,
	}
	for _, opt := range opts {
		opt(&wr)
	}
	return wr
}

func (w Writer) Close() error {
	if c, ok := w.w.(io.WriteCloser); ok {
		return c.Close()
	}
	return fmt.Errorf("given writer is no WriteCloser")
}

func (w Writer) Write(data interface{}) error {
	j, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("could not json marshal data: %w", err)
	}

	if w.prefix != "" {
		_, err = w.w.Write([]byte(w.prefix))
		if err != nil {
			return fmt.Errorf("could not write prefix to underlying io.Writer: %w", err)
		}
	}

	_, err = w.w.Write(j)
	if err != nil {
		return fmt.Errorf("could not write json data to underlying io.Writer: %w", err)
	}

	_, err = w.w.Write([]byte("\n"))
	if err != nil {
		return fmt.Errorf("could not write newline to underlying io.Writer: %w", err)
	}

	if f, ok := w.w.(http.Flusher); ok {
		// If http writer, flush as well
		f.Flush()
	}
	return nil
}
