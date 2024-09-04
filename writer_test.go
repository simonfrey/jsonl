package jsonl_test

import (
	"bytes"
	"github.com/simonfrey/jsonl"
	"testing"
)

type T1 struct {
	Type  string `json:"type"`
	Count int    `json:"am"`
}
type T2 struct {
	Type    string `json:"type"`
	Comment string `json:"am"`
}

func TestWriter_Write(t *testing.T) {
	expectedResponse := "{\"type\":\"T1\",\"am\":2}\n"

	buff := bytes.Buffer{}

	w := jsonl.NewWriter(&buff)

	err := w.Write(T1{
		Type:  "T1",
		Count: 2,
	})
	if err != nil {
		t.Fatalf("could not write T1: %s", err)
	}

	if buff.String() != expectedResponse {
		t.Fatalf("Response %q is not as expected %q", buff.String(), expectedResponse)
	}
}

func TestWriter_WriteMultipleLines(t *testing.T) {
	expectedResponse := "{\"type\":\"T1\",\"am\":2}\n{\"type\":\"T2\",\"am\":\"I am T2\"}\n{\"type\":\"T1\",\"am\":9999}\n"

	buff := bytes.Buffer{}

	w := jsonl.NewWriter(&buff)

	err := w.Write(T1{
		Type:  "T1",
		Count: 2,
	})
	if err != nil {
		t.Fatalf("could not write T1: %s", err)
	}
	err = w.Write(T2{
		Type:    "T2",
		Comment: "I am T2",
	})
	if err != nil {
		t.Fatalf("could not write T2: %s", err)
	}
	err = w.Write(T1{
		Type:  "T1",
		Count: 9999,
	})
	if err != nil {
		t.Fatalf("could not write T1: %s", err)
	}

	if buff.String() != expectedResponse {
		t.Fatalf("Response %q is not as expected %q", buff.String(), expectedResponse)
	}
}
