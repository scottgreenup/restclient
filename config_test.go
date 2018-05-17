package restclient

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewConfig(t *testing.T) {
	config := NewRequestBuilder()
	require.Error(t, config.Validate())
}

func TestWithEndpoint(t *testing.T) {
	config := NewRequestBuilder().WithEndpoint("https://api.trello.com/1/")
	err := config.Validate()
	require.NoError(t, err)
}

func TestNewMutator(t *testing.T) {
	rm := NewRequestMutator()
	require.Len(t, rm.mutations, 0)
}

func TestURL(t *testing.T) {
	rm := NewRequestMutator(
		BaseURL("https://scottgreenup.com/"),
		ResolvePath("/api/whatever"),
	)

	req, err := rm.NewRequest()
	require.NoError(t, err)
	require.Equal(t, req.URL.String(), "https://scottgreenup.com/api/whatever")
}
