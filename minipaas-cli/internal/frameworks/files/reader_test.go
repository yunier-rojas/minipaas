package files

import (
	"os"
	"path/filepath"
	"testing"
)

func TestReadContent_FileWithoutNameUsesBaseName(t *testing.T) {
	dir := t.TempDir()
	file := filepath.Join(dir, "app.conf")
	if err := os.WriteFile(file, []byte("file-content"), 0644); err != nil {
		t.Fatal(err)
	}

	baseName, content, err := NewReaderService().ReadContent(file, "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if baseName != "app.conf" {
		t.Fatalf("baseName = %q, want app.conf", baseName)
	}
	if string(content) != "file-content" {
		t.Fatalf("content = %q", content)
	}
}

func TestReadContent_FileWithNameNameWins(t *testing.T) {
	dir := t.TempDir()
	file := filepath.Join(dir, "app.conf")
	if err := os.WriteFile(file, []byte("file-content"), 0644); err != nil {
		t.Fatal(err)
	}

	baseName, content, err := NewReaderService().ReadContent(file, "custom")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if baseName != "custom" {
		t.Fatalf("baseName = %q, want custom", baseName)
	}
	if string(content) != "file-content" {
		t.Fatalf("content = %q", content)
	}
}

func TestReadContent_StdinWithName(t *testing.T) {
	reader, writer, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	original := os.Stdin
	os.Stdin = reader
	defer func() { os.Stdin = original }()

	go func() {
		_, _ = writer.WriteString("stdin-content")
		_ = writer.Close()
	}()

	baseName, content, err := NewReaderService().ReadContent("", "postgres_password")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if baseName != "postgres_password" {
		t.Fatalf("baseName = %q, want postgres_password", baseName)
	}
	if string(content) != "stdin-content" {
		t.Fatalf("content = %q", content)
	}
}

func TestReadContent_StdinWithoutNameErrors(t *testing.T) {
	if _, _, err := NewReaderService().ReadContent("", ""); err == nil {
		t.Fatalf("expected error when neither file nor name is provided")
	}
}

func TestReadContent_MissingFileErrors(t *testing.T) {
	if _, _, err := NewReaderService().ReadContent(filepath.Join(t.TempDir(), "nope"), "x"); err == nil {
		t.Fatalf("expected error for missing file")
	}
}
