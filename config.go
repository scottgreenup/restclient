package restclient

import (
	"net/http"
	"net/url"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/pkg/errors"
)

// Config represents the common ground that the endpoints of a RESTful API has. Attributes like the HTTP client, and the
// base URL
type Config struct {
	Endpoint   *url.URL
	HTTPClient *http.Client

	merr *multierror.Error
}

// NewConfig generates a new configuration
func NewConfig() *Config {
	return &Config{
		HTTPClient: http.DefaultClient,
	}
}

// NewOperation generates a new operation from the Config, giving the operation the common ground set by the config
func (c *Config) NewOperation(method string) (*Operation, error) {
	if err := c.Validate(); err != nil {
		return nil, err
	}

	return newOperation(method, c), nil
}

// WithEndpoint sets the base URL, i.e. the URL prefix for all of our API calls
func (c *Config) WithEndpoint(u string) *Config {
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

// WithHTTPClient sets a custom http client to use, by default the http.DefaultClient is used
func (c *Config) WithHTTPClient(client *http.Client) *Config {
	if client == nil {
		c.merr = multierror.Append(c.merr, errors.New("WithHTTPClient called with nil"))
		return c
	}

	c.HTTPClient = client
	return c
}

// Validate will validate the config
func (c *Config) Validate() error {
	var merr *multierror.Error

	if c.merr == nil {
		merr = multierror.Append(merr, c.merr)
	}

	if c.Endpoint == nil {
		merr = multierror.Append(merr, errors.New("endpoint is nil"))
	}

	return merr.ErrorOrNil()
}
