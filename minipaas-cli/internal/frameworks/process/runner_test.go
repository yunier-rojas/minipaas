package process

import (
	"os"
	"path/filepath"
	"testing"
)

func makeScript(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, "script.sh")
	if err := os.WriteFile(p, []byte("#!/bin/sh\n"+content+"\n"), 0755); err != nil {
		t.Fatalf("write script: %v", err)
	}
	return p
}

func TestRunner_RunSuccessAndFailure(t *testing.T) {
	runner := NewRunner()

	ok := makeScript(t, "exit 0")
	if err := runner.Run([]string{ok}, false); err != nil {
		t.Fatalf("expected success, got: %v", err)
	}

	fail := makeScript(t, "exit 1")
	if err := runner.Run([]string{fail}, false); err == nil {
		t.Fatalf("expected failure, got nil")
	}
}

func TestRunner_RunWithInput(t *testing.T) {
	runner := NewRunner()
	sc := makeScript(t, `
in=$(cat)
if [ "$in" = "hello" ]; then
  exit 0
else
  exit 1
fi`)

	if err := runner.RunWithInput([]string{sc}, []byte("hello"), false); err != nil {
		t.Fatalf("expected success, got: %v", err)
	}
	if err := runner.RunWithInput([]string{sc}, []byte("nope"), false); err == nil {
		t.Fatalf("expected failure for wrong input")
	}
}

func TestRunner_RunOutput(t *testing.T) {
	runner := NewRunner()
	sc := makeScript(t, `printf 'hello\n'`)

	out, err := runner.RunOutput([]string{sc}, false)
	if err != nil {
		t.Fatalf("expected success, got: %v", err)
	}
	if out != "hello\n" {
		t.Fatalf("output = %q, want %q", out, "hello\n")
	}
}
