package restclient

import (
	"encoding/base64"
	"testing"

	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
)

func TestJSON(t *testing.T) {
	type Foo struct {
		A int `json:"element_a"`
	}

	x := Foo{A: 1}

	config := NewRequestBuilder().
		WithEndpoint("https://api.trello.com/1/")

	operation, err := config.NewOperation(http.MethodGet)
	require.NoError(t, err)

	_, err = operation.
		BodyFromJSON(x).
		WithPath("/boards/{id}").
		WithPathVar("id", "hello").
		BuildRequest()

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

	config := NewRequestBuilder().WithEndpoint("https://api.trello.com/1/")

	for _, tc := range cases {
		operation, err := config.NewOperation(http.MethodGet)
		require.NoError(t, err)

		_, err = operation.WithPath("endpoint").BodyFromJSONString(tc.Input).BuildRequest()
		require.NoError(t, err)
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

	config := NewRequestBuilder().WithEndpoint("https://api.trello.com/1/")

	for _, tc := range cases {
		operation, err := config.NewOperation(http.MethodGet)
		require.NoError(t, err)

		_, err = operation.WithPath("endpoint").BodyFromJSONString(tc.Input).BuildRequest()
		require.Error(t, err)
	}
}

func TestFullExample(t *testing.T) {

	type AuthenticationStruct struct {
		Username string
		Password string
	}

	BasicAuthenticationMethod := func(req *http.Request, v ...interface{}) (*http.Request, error) {
		auth := v[0].(AuthenticationStruct)
		code := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", auth.Username, auth.Password)))
		req.Header.Add("Authorization", fmt.Sprintf("Basic %s", code))
		return req, nil
	}

	config := NewRequestBuilder().
		WithEndpoint("https://api.trello.com/1/").
		WithAuthenticationMethod(BasicAuthenticationMethod)

	operation, _ := config.NewOperation(http.MethodPost)

	req, err := operation.
		WithPath("boards/{boardid}/cards/{cardid}").
		WithPathVar("boardid", "12345").
		WithPathVar("cardid", "12345").
		BodyFromJSONString(`{"name": "Seymour Butts"}`).
		Authenticate(AuthenticationStruct{Username: "", Password: ""}).
		BuildRequest()

	require.NoError(t, err)

	//resp, err := http.DefaultClient.Do(req)

	require.Equal(t, "https://api.trello.com/1/boards/12345/cards/12345", req.URL.String())
	require.Equal(t, "application/json", req.Header.Get("Content-Type"))

}

func TestWithPathValid(t *testing.T) {
	cases := []struct {
		Path     string
		PathVars map[string]string
	}{
		{"/boards/", map[string]string{}},
		{"/boards/{id}", map[string]string{"id": "10"}},
		{"/boards/{id}/{cardid}", map[string]string{"id": "10", "cardid": "123"}},
		{"/boards/{id}-{cardid}", map[string]string{"id": "10", "cardid": "123"}},
		{"/boards/{id}/cards/{cardid}", map[string]string{"id": "10", "cardid": "123"}},
	}

	config := NewRequestBuilder().WithEndpoint("https://api.trello.com/1/")

	for _, tc := range cases {
		operation, err := config.NewOperation(http.MethodGet)
		require.NoError(t, err)

		operation = operation.WithPath(tc.Path)

		for k, v := range tc.PathVars {
			operation.WithPathVar(k, v)
		}

		_, err = operation.BuildRequest()
		require.NoErrorf(t, err, "WithPath(%s).WithPathVars(%v)", tc.Path, tc.PathVars)
	}
}

func TestWithPathInvalid(t *testing.T) {
	cases := []struct {
		Path     string
		PathVars map[string]string
	}{
		{"/boards/{id}", map[string]string{"ID": "10"}},
		{"/boards/{id}", map[string]string{}},
		{"/boards/", map[string]string{"id": "10"}},
		{"/boards/{hello}", map[string]string{"id": "10"}},
		{"/boards/{hello}", map[string]string{"hello": "20", "id": "10"}},
	}

	config := NewRequestBuilder().WithEndpoint("https://api.trello.com/1/")

	for _, tc := range cases {
		operation, err := config.NewOperation(http.MethodGet)
		require.NoError(t, err)

		operation = operation.WithPath(tc.Path)

		for k, v := range tc.PathVars {
			operation.WithPathVar(k, v)
		}

		_, err = operation.BuildRequest()
		require.Errorf(t, err, "WithPath(%s).WithPathVars(%v)", tc.Path, tc.PathVars)
	}
}

func TestFormatValid(t *testing.T) {
	cases := []struct {
		Format        string
		Vars          map[string]string
		ExpectedValue string
	}{
		{"nothing special", nil, "nothing special"},
		{"a", nil, "a"},
		{"!@#$%^&*()", nil, "!@#$%^&*()"},
		{"{id}", map[string]string{"id": "ACDC"}, "ACDC"},
		{"{id}{id}", map[string]string{"id": "Titanium"}, "TitaniumTitanium"},
		{"{id}/{id}", map[string]string{"id": "Titanium"}, "Titanium/Titanium"},
		{"{first_name}/{last_name}", map[string]string{"first_name": "Bruce", "last_name": "Wayne"}, "Bruce/Wayne"},
		{"{recursion_sucks}", map[string]string{"recursion_sucks": "{recursion_does_not_work}"}, "{recursion_does_not_work}"},
	}

	for _, tc := range cases {
		val, err := format(tc.Format, tc.Vars)
		assert.Equal(t, tc.ExpectedValue, val)
		assert.NoError(t, err)
	}
}

func TestFormatInvalid(t *testing.T) {
	cases := []struct {
		Format string
		Vars   map[string]string
	}{
		{"value/{id}", nil},
	}

	for _, tc := range cases {
		val, err := format(tc.Format, tc.Vars)
		assert.Equal(t, "", val)
		assert.Error(t, err)
	}
}
