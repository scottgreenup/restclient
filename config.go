package restclient

import (
    "net/http"
    "net/url"
    "strings"

    "github.com/hashicorp/go-multierror"
    "github.com/pkg/errors"
)

type Config struct {
    endpoint   *url.URL
    httpClient *http.Client

    merr *multierror.Error
}

func NewConfig() *Config {
    return &Config{
        httpClient: http.DefaultClient,
    }
}

func (c *Config) WithEndpoint(u string) (*Config) {
    if !strings.HasPrefix(u, "http://") && !strings.HasPrefix(u, "https://") {
        c.merr = multierror.Append(c.merr, errors.New("WithEndpoint endpoint is missing http:// or https:// prefix"))
        return c
    }

    endpoint, err := url.Parse(u)
    if err != nil {
        c.merr = multierror.Append(c.merr, errors.Wrap(err, "WithEndpoint endpoint could not be parsed"))
        return c
    }

    c.endpoint = endpoint
	return c
}

func (c *Config) WithHTTPClient(client *http.Client) (*Config) {
    if client == nil {
        c.merr = multierror.Append(c.merr, errors.New("WithHTTPClient called with nil"))
        return c
    }

    c.httpClient = client
    return c
}

func (c *Config) Validate() (err error) {
    if c.merr == nil {
        err = multierror.Append(err, c.merr)
    }

    if c.endpoint == nil {
        err = multierror.Append(err, errors.New("endpoint is nil"))
    }

    return
}
