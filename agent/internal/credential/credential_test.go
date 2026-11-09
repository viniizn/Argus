package credential

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

func TestLoadNotFound(t *testing.T) {
	_, err := Load(filepath.Join(t.TempDir(), "missing.json"))
	if !errors.Is(err, ErrNotFound) {
		t.Fatalf("want ErrNotFound, got %v", err)
	}
}

func TestSaveThenLoad(t *testing.T) {
	path := filepath.Join(t.TempDir(), "sub", "credential.json")

	if err := Save(path, Secret("abc123")); err != nil {
		t.Fatalf("Save: %v", err)
	}
	got, err := Load(path)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if got.Reveal() != "abc123" {
		t.Fatalf("got %q, want abc123", got.Reveal())
	}
}

func TestSavePermissions(t *testing.T) {
	path := filepath.Join(t.TempDir(), "agent", "credential.json")
	if err := Save(path, Secret("abc123")); err != nil {
		t.Fatalf("Save: %v", err)
	}

	info, err := os.Stat(path)
	if err != nil {
		t.Fatal(err)
	}
	if perm := info.Mode().Perm(); perm != 0o600 {
		t.Fatalf("file perm = %o, want 600", perm)
	}
	dirInfo, err := os.Stat(filepath.Dir(path))
	if err != nil {
		t.Fatal(err)
	}
	if perm := dirInfo.Mode().Perm(); perm != 0o700 {
		t.Fatalf("dir perm = %o, want 700", perm)
	}
}

func TestSaveOverwritesAtomicallyAndLeavesNoTemp(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "credential.json")

	if err := Save(path, Secret("first")); err != nil {
		t.Fatal(err)
	}
	if err := Save(path, Secret("second")); err != nil {
		t.Fatal(err)
	}
	got, _ := Load(path)
	if got.Reveal() != "second" {
		t.Fatalf("got %q, want second", got.Reveal())
	}

	entries, _ := os.ReadDir(dir)
	if len(entries) != 1 {
		t.Fatalf("expected only the credential file, found %d entries", len(entries))
	}
}

func TestLoadRejectsBadContent(t *testing.T) {
	tests := []struct {
		name    string
		content string
	}{
		{"invalid json", "not json"},
		{"empty credential", `{"credential":""}`},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path := filepath.Join(t.TempDir(), "credential.json")
			if err := os.WriteFile(path, []byte(tt.content), 0o600); err != nil {
				t.Fatal(err)
			}
			if _, err := Load(path); err == nil {
				t.Fatal("expected error, got nil")
			}
		})
	}
}

func TestSecretIsRedactedWhenFormatted(t *testing.T) {
	s := Secret("super-secret-value")
	for _, out := range []string{
		fmt.Sprint(s), fmt.Sprintf("%v", s), fmt.Sprintf("%s", s), fmt.Sprintf("%#v", s),
	} {
		if out != "<redacted>" && out != "credential.Secret(<redacted>)" {
			t.Fatalf("secret leaked in formatting: %q", out)
		}
	}
}