package restclient

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
)

func TestJSON(t *testing.T) {

	type Foo struct {
		A int `json:"element_a"`
	}

	x := Foo{A: 1}

	config := NewConfig().
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
	panic("not yet implemented")
}

func TestJSONStringInvalid(t *testing.T) {
	panic("not yet implemented")
}

func TestFullExample(t *testing.T) {
	panic("not yet implemented")
}

func TestWithPathValid(t *testing.T) {
	panic("not yet implemented")
}

func TestWithPathInvalid(t *testing.T) {
	panic("not yet implemented")
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
