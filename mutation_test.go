package restclient

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestURL(t *testing.T) {
	rm := NewRequestMutator(
		BaseURL("https://scottgreenup.com/"),
		ResolvePath("/api/whatever"),
	)

	req, err := rm.NewRequest()
	require.NoError(t, err)
	require.Equal(t, req.URL.String(), "https://scottgreenup.com/api/whatever")
}

func TestHeader(t *testing.T) {
	rm := NewRequestMutator()
	req, err := rm.NewRequest(
		Header("abc", "123"),
	)
	require.NoError(t, err)
	require.Equal(t, req.Header.Get("abc"), "123")
}

func TestJSON(t *testing.T) {
	type Foo struct {
		A int `json:"element_a"`
	}

	data := &Foo{A: 1}

	handler := func(w http.ResponseWriter, req *http.Request) {
		defer req.Body.Close() // nolint: errcheck

		target := &Foo{}
		decoder := json.NewDecoder(req.Body)
		err := decoder.Decode(target)
		require.NoError(t, err)
		require.Equal(t, target.A, 1)
	}
	srv := httptest.NewServer(http.HandlerFunc(handler))
	defer srv.Close()

	req, err := NewRequestMutator(
		BaseURL(srv.URL),
	).NewRequest(
		Method(http.MethodGet),
		BodyFromJSON(data),
	)

	require.NoError(t, err)
	require.NotNil(t, req.Body)
	require.NotNil(t, req.GetBody)

	// +1 for newline
	require.Equal(t, req.ContentLength, int64(len(`{"element_a":1}`)+1))

	// Activate our handler
	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, resp.Body.Close())
	require.NoError(t, err)
}

func TestJSONStringValid(t *testing.T) {
	cases := []struct {
		Input string
	}{
		{`{}`},
		{`{"words": 1}`},
		{`{"words": 1, "child": {"a": "words"}}`},
	}

	rm := NewRequestMutator(
		BaseURL("https://api.trello.com/1/"),
	)

	for _, tc := range cases {
		req, err := rm.NewRequest(
			Method(http.MethodGet),
			BodyFromJSONString(tc.Input),
		)
		require.NoError(t, err)
		require.NotNil(t, req.Body)
		require.NotNil(t, req.GetBody)
		require.Equal(t, req.ContentLength, int64(len(tc.Input)))
	}
}

func TestJSONStringInvalid(t *testing.T) {
	cases := []struct {
		Input string
	}{
		{`{`},
		{`}`},
		{`{"words": 1", "child": {"a": "words"}}`},
		{`content=3`},
	}

	rm := NewRequestMutator(
		BaseURL("https://api.trello.com/1/"),
	)

	for _, tc := range cases {
		_, err := rm.NewRequest(
			Method(http.MethodGet),
			BodyFromJSONString(tc.Input),
		)
		require.Error(t, err)
	}
}
