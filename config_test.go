package restclient

import (
    "testing"

    "github.com/stretchify/require"
)

func TestNewConfig(t *testing.T) {
    config := NewConfig()
    err := config.Validate()

    require.Error(t, err)
}
