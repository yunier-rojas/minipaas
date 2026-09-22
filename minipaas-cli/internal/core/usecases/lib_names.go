package usecases

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
)

func hashedName(baseName string, content []byte) string {
	hash := sha256.Sum256(content)
	prefix := hex.EncodeToString(hash[:])[:8]
	return fmt.Sprintf("%s.%s", baseName, prefix)
}

func secretHashName(baseName string, content []byte) string {
	return hashedName(baseName, content)
}

func configHashName(baseName string, content []byte) string {
	return hashedName(baseName, content)
}
