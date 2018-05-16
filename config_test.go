package restclient

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewConfig(t *testing.T) {
	config := NewConfig()
	require.Error(t, config.Validate())
}

func TestWithEndpoint(t *testing.T) {
	config := NewConfig().WithEndpoint("https://api.trello.com/1/")
	err := config.Validate()
	require.NoError(t, err)
}
