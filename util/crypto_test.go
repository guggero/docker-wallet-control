package util

import (
    "testing"
    "github.com/stretchr/testify/assert"
)

func TestHashPassword(t *testing.T) {
    // given
    password := "test"
    salt := "verysalty"

    // when
    result := HashPassword(password, salt)

    // then
    assert.Equal(t, "c32d9b1ea5e69ef67aebe6d3f25a256ffa1c58e9b7ecae072cd50916c6674ad5", result)
}
