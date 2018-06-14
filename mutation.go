package restclient

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/pkg/errors"
)

type RequestMutation func(req *http.Request) (*http.Request, error)

const (
	// ContentTypeJSON is the default for JSON strings
	ContentTypeJSON = "application/json"
)

func Body(body io.Reader) RequestMutation {
	return func(req *http.Request) (*http.Request, error) {
		if body == nil {
			req.Body = nil
			return req, nil
		}

		// If the ContentLength is actually 0, we can signal that ContentLength is ACTUALLY 0, and not just unknown
		OptimizeIfEmpty := func(req *http.Request) {
			if req.ContentLength == 0 {

				// This signals that the ContentLengt
				req.Body = http.NoBody
				req.GetBody = func() (io.ReadCloser, error) {
					return http.NoBody, nil
				}
			}
		}

		switch v := body.(type) {
		case *bytes.Buffer:
			req.ContentLength = int64(v.Len())
			buf := v.Bytes()
			req.Body = ioutil.NopCloser(body)
			req.GetBody = func() (io.ReadCloser, error) {
				r := bytes.NewReader(buf)
				return ioutil.NopCloser(r), nil
			}
			OptimizeIfEmpty(req)

		case *bytes.Reader:
			req.ContentLength = int64(v.Len())
			snapshot := *v
			req.Body = ioutil.NopCloser(body)
			req.GetBody = func() (io.ReadCloser, error) {
				r := snapshot
				return ioutil.NopCloser(&r), nil
			}
			OptimizeIfEmpty(req)

		case *strings.Reader:
			req.ContentLength = int64(v.Len())
			snapshot := *v
			req.Body = ioutil.NopCloser(body)
			req.GetBody = func() (io.ReadCloser, error) {
				r := snapshot
				return ioutil.NopCloser(&r), nil
			}
			OptimizeIfEmpty(req)

		default:
			// We don't have the same backwards compatibility issues as req.go:
			// https://github.com/golang/go/blob/226651a541/src/net/http/req.go#L823-L839

			// Convert the body to a ReadCloser
			rc, ok := body.(io.ReadCloser)
			if !ok {
				rc = ioutil.NopCloser(body)
			}
			req.Body = rc
			req.GetBody = func() (io.ReadCloser, error) {
				return rc, nil
			}

			// We don't know how to get the Length, this is signalled by setting ContentLength to 0
			req.ContentLength = 0
		}

		return req, nil
	}
}

func Context(ctx context.Context) RequestMutation {
	return func(req *http.Request) (*http.Request, error) {
		return req.WithContext(ctx), nil
	}
}

// BodyFromJSON marshals the interface `v` into JSON for the req
func BodyFromJSON(v interface{}) RequestMutation {
	return func(req *http.Request) (*http.Request, error) {
		b := new(bytes.Buffer)
		if err := json.NewEncoder(b).Encode(v); err != nil {
			return nil, errors.Wrap(err, "could not encode body")
		}
		req, err := Body(b)(req)
		if err != nil {
			return nil, errors.Wrap(err, "could not set body")
		}
		req.Header.Set("Content-Type", ContentTypeJSON)
		return req, nil
	}
}

// BodyFromJSONString uses the string as JSON for the req
func BodyFromJSONString(s string) RequestMutation {
	return func(req *http.Request) (*http.Request, error) {
		if err := json.Unmarshal([]byte(s), &map[string]interface{}{}); err != nil {
			return nil, errors.Wrap(err, "BodyFromJSONString received string that was not valid JSON")
		}
		req, err := Body(strings.NewReader(s))(req)
		if err != nil {
			return nil, errors.Wrap(err, "could not set body")
		}
		req.Header.Set("Content-Type", ContentTypeJSON)
		return req, nil
	}
}

func Method(method string) RequestMutation {
	return func(req *http.Request) (*http.Request, error) {
		req.Method = method
		return req, nil
	}
}

// Header adds a header to the request
func Header(key, value string) RequestMutation {
	return func(req *http.Request) (*http.Request, error) {
		req.Header.Add(key, value)
		return req, nil
	}
}

// BaseURL sets the URL of the req from a URL string
func BaseURL(base string) RequestMutation {
	return func(req *http.Request) (*http.Request, error) {
		u, err := url.Parse(base)
		if err != nil {
			return nil, err
		}
		req.URL = u
		return req, nil
	}
}

// ResolvePath sets the path of the req by resolving it on the req.URL
func ResolvePath(path string) RequestMutation {
	return func(req *http.Request) (*http.Request, error) {
		u, err := url.Parse(path)
		if err != nil {
			return nil, err
		}
		req.URL.Path = req.URL.ResolveReference(u).Path
		return req, nil
	}
}
