package util

import (
    "crypto/sha256"
    "fmt"
)

func HashPassword(password string, salt string) (string) {
    hashBytes := sha256.Sum256([]byte(salt + password))
    return fmt.Sprintf("%x", hashBytes)
}