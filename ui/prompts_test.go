package ui

import (
	"bufio"
	"bytes"
	"context"
	"testing"
)

func TestBuildPrompt(t *testing.T) {
	table := []struct {
		Input          string
		ExpectedOutput string
		ExpectedResult bool
	}{
		{
			Input:          "true\n",
			ExpectedOutput: "Build a new GoBlog binary? ('true', 'false')\n> ",
			ExpectedResult: true,
		},
		{
			Input:          "false\n",
			ExpectedOutput: "Build a new GoBlog binary? ('true', 'false')\n> ",
			ExpectedResult: false,
		},
	}
	for _, tt := range table {
		t.Run(tt.Input, func(t *testing.T) {
			in, out := bytes.Buffer{}, bytes.Buffer{}
			inBuffered, outBuffered := bufio.NewReader(&in), bufio.NewWriter(&out)
			rw := bufio.ReadWriter{
				Reader: inBuffered,
				Writer: outBuffered,
			}

			// inject the input, itll just be
			// waiting for the first read.
			in.WriteString(tt.Input)

			p := Prompter{
				stdio: rw,
			}

			b, err := p.ShouldBuild(context.Background())
			if err != nil {
				t.Fatalf("%v", err)
			}
			if b != tt.ExpectedResult {
				t.Logf("got: %v, want: %v", b, tt.ExpectedResult)
				t.Fail()
			}

			// confirm prompt is correct
			if out.String() != tt.ExpectedOutput {
				t.Logf("got: %v, want: %v", out.String(), tt.ExpectedResult)
				t.Fail()
			}
		})
	}
}

func TestPublishPostPrompt(t *testing.T) {
	table := []struct {
		Input          string
		ExpectedOutput string
		ExpectedResult bool
	}{
		{
			Input:          "true\n",
			ExpectedOutput: "Publish this post? ('true', 'false')\n> ",
			ExpectedResult: true,
		},
		{
			Input:          "false\n",
			ExpectedOutput: "Publish this post? ('true', 'false')\n> ",
			ExpectedResult: false,
		},
	}
	for _, tt := range table {
		t.Run(tt.Input, func(t *testing.T) {
			in, out := bytes.Buffer{}, bytes.Buffer{}
			inBuffered, outBuffered := bufio.NewReader(&in), bufio.NewWriter(&out)
			rw := bufio.ReadWriter{
				Reader: inBuffered,
				Writer: outBuffered,
			}

			// inject the input, itll just be
			// waiting for the first read.
			in.WriteString(tt.Input)

			p := Prompter{
				stdio: rw,
			}

			b, err := p.ShouldPublishPost(context.Background())
			if err != nil {
				t.Fatalf("%v", err)
			}
			if b != tt.ExpectedResult {
				t.Logf("got: %v, want: %v", b, tt.ExpectedResult)
				t.Fail()
			}

			// confirm prompt is correct
			if out.String() != tt.ExpectedOutput {
				t.Logf("got: %v, want: %v", out.String(), tt.ExpectedResult)
				t.Fail()
			}
		})
	}
}
