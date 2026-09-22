package usecases

import (
	"crypto/sha256"
	"encoding/hex"
	"testing"
)

func TestHashedName(t *testing.T) {
	content := []byte("hello")
	hash := sha256.Sum256(content)
	want := "app.conf." + hex.EncodeToString(hash[:])[:8]

	if got := hashedName("app.conf", content); got != want {
		t.Fatalf("hashedName() = %q, want %q", got, want)
	}
}

func TestHashNameFunctions(t *testing.T) {
	tests := []struct {
		name     string
		baseName string
		content  []byte
		build    func(string, []byte) string
	}{
		{name: "secret", baseName: "postgres_password", content: []byte("hunter2"), build: secretHashName},
		{name: "config", baseName: "app.json", content: []byte("{\"a\":1}"), build: configHashName},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hash := sha256.Sum256(tt.content)
			want := tt.baseName + "." + hex.EncodeToString(hash[:])[:8]

			if got := tt.build(tt.baseName, tt.content); got != want {
				t.Fatalf("hash name = %q, want %q", got, want)
			}
		})
	}
}

func TestSecretHashNameSamePrefixAsGeneric(t *testing.T) {
	if got, want := secretHashName("a", []byte("b")), hashedName("a", []byte("b")); got != want {
		t.Fatalf("secretHashName() = %q, want %q", got, want)
	}
	if got, want := configHashName("a", []byte("b")), hashedName("a", []byte("b")); got != want {
		t.Fatalf("configHashName() = %q, want %q", got, want)
	}
}
