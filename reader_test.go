package jsonl_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/simonfrey/jsonl"
	"strings"
	"testing"
)

func TestReader_ReadSingleLine(t *testing.T) {
	input := "{\"type\":\"T1\",\"am\":2}\n"

	r := jsonl.NewReader(strings.NewReader(input))

	t1 := T1{}
	err := r.ReadSingleLine(&t1)
	if err != nil {
		t.Fatalf("could not ReadSingleLine: %s", err)
	}

	if t1.Count != 2 {
		t.Fatalf("read invalid count. Got %d but exepcted %d", t1.Count, 2)
	}
}

func TestReader_ReadLines(t *testing.T) {
	input := "{\"type\":\"T1\",\"am\":2}\n{\"type\":\"T2\",\"am\":\"I am T2\"}\n{\"type\":\"T1\",\"am\":9999}\n"
	r := jsonl.NewReader(strings.NewReader(input))

	output := ""
	err := r.ReadLines(func(data []byte) error {
		switch {
		case bytes.Contains(data, []byte(`"T1"`)):
			// T1 struct type
			t := T1{}
			err := json.Unmarshal(data, &t)
			if err != nil {
				return fmt.Errorf("could not unmarshal into T1: %w", err)
			}
			output += fmt.Sprintf("%T:%d|", t, t.Count)
		case bytes.Contains(data, []byte(`"T2"`)):
			// T2 struct type
			t := T2{}
			err := json.Unmarshal(data, &t)
			if err != nil {
				return fmt.Errorf("could not unmarshal into T2: %w", err)
			}
			output += fmt.Sprintf("%T:%s|", t, t.Comment)
		}
		return nil
	})
	if err != nil {
		t.Fatalf("could not read lines: %s", err)
	}

	if output != "jsonl_test.T1:2|jsonl_test.T2:I am T2|jsonl_test.T1:9999|" {
		t.Fatalf("did get wrong response: %q", output)
	}
}
