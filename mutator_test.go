package restclient

import (
	"testing"

	"encoding/base64"
	"fmt"
	"github.com/stretchr/testify/require"
	"net/http"
)

func TestNewMutator(t *testing.T) {
	rm := NewRequestMutator()
	require.Len(t, rm.mutations, 0)
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

	req, _ := rm.NewRequest(
		ResolvePath("/api/whatever"),
		BasicAuthMutator("admin", "supersecretpassword"),
	)

	fmt.Println("URL: ", req.URL.String())
	fmt.Println("Authorization Header: ", req.Header.Get("Authorization"))
}
