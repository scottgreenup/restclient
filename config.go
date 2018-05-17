package restclient

import (
	"net/http"
	"net/url"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/pkg/errors"
)

// Config represents the common ground that the endpoints of a RESTful API has. Attributes like the HTTP client, and the
// base URL. This acts as an OperationBuilder, where the OperationBuilder acts as a RequestBuilder.
type RequestBuilder struct {
	Endpoint   *url.URL
	AuthMethod AuthenticationMethod

	merr *multierror.Error
}

// NewConfig generates a new configuration
func NewRequestBuilder() *RequestBuilder {
	return &RequestBuilder{}
}

// AuthenticationMethod represents the user-made way of setting a request to be authenticated
type AuthenticationMethod func(req *http.Request, v ...interface{}) (*http.Request, error)

// SetAuthenticationMethod sets the AuthenticationMethod for operations, which is invoked with Operation.Authenticate()
func (c *RequestBuilder) WithAuthenticationMethod(authMethod AuthenticationMethod) *RequestBuilder {
	c.AuthMethod = authMethod
	return c
}

// NewOperation generates a new operation from the Config, giving the operation the common ground set by the config
func (c *RequestBuilder) NewOperation(method string) (*Request, error) {
	if err := c.Validate(); err != nil {
		return nil, err
	}

	return newRequest(method, c), nil
}

// WithEndpoint sets the base URL, i.e. the URL prefix for all of our API calls
func (c *RequestBuilder) WithEndpoint(u string) *RequestBuilder {
	if !strings.HasPrefix(u, "http://") && !strings.HasPrefix(u, "https://") {
		c.merr = multierror.Append(c.merr, errors.New("WithEndpoint endpoint is missing http:// or https:// prefix"))
		return c
	}

	endpoint, err := url.Parse(u)
	if err != nil {
		c.merr = multierror.Append(c.merr, errors.Wrap(err, "WithEndpoint endpoint could not be parsed"))
		return c
	}

	c.Endpoint = endpoint
	return c
}

// Validate will validate the config
func (c *RequestBuilder) Validate() error {
	var merr *multierror.Error

	if c.merr == nil {
		merr = multierror.Append(merr, c.merr)
	}

	if c.Endpoint == nil {
		merr = multierror.Append(merr, errors.New("endpoint is nil"))
	}

	return merr.ErrorOrNil()
}
