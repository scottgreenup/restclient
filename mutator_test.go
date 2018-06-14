package restclient

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewMutator(t *testing.T) {
	rm := NewRequestMutator()
	require.Len(t, rm.mutations, 0)
}

func ExampleRequestMutator() {

	type ExampleRequest struct {
		Version uint64 `json:"number"`
		Hash    string `json:"hash"`
	}

	requestData := &ExampleRequest{
		Version: 1,
		Hash:    "1c76f1f33f12c14a63026f71c8d17ab2",
	}

	rm := NewRequestMutator(
		BaseURL("https://scottgreenup.com/"),
	)

	BasicAuthMutator := func(username, password string) RequestMutation {
		return func(req *http.Request) (*http.Request, error) {
			code := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", username, password)))
			req.Header.Add("Authorization", fmt.Sprintf("Basic %s", code))
			return req, nil
		}
	}

	req, _ := rm.NewRequest(
		ResolvePath("/api/example"),
		Method(http.MethodPost),
		BasicAuthMutator("MyUsername", "4RuwRmDkLm990qkXMK6obWK88S7pW3K3"),
		BodyFromJSON(requestData),
	)

	fmt.Println("URL:", req.URL.String())
	fmt.Println("Content Length:", req.ContentLength)
	fmt.Println("Authorization Header:", req.Header.Get("Authorization"))

	// output:
	// URL: https://scottgreenup.com/api/example
	// Content Length: 55
	// Authorization Header: Basic TXlVc2VybmFtZTo0UnV3Um1Ea0xtOTkwcWtYTUs2b2JXSzg4UzdwVzNLMw==
}
