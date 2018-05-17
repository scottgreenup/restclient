package restclient

import (
	"testing"

	"encoding/base64"
	"fmt"
	"github.com/stretchr/testify/require"
	"net/http"
	"strings"
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

func ExampleAuthenticationMethod() {
	rm := NewRequestMutator(
		BaseURL("https://scottgreenup.com/"),
	)

	BasicAuthMutator := func(username, password string) RequestMutation {
		return func(req *http.Request) error {
			code := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", username, password)))
			req.Header.Add("Authorization", fmt.Sprintf("Basic %s", code))
			return nil
		}
	}

	req, err := rm.NewRequest(
		ResolvePath("/api/whatever"),
		BasicAuthMutator("admin", "supersecretpassword"),
	)

	fmt.Println("URL: ", req.URL.String())
	fmt.Println("Authorization Header: ", req.Header.Get("Authorization"))
}
