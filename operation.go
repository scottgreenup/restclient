package restclient

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/pkg/errors"
)

const (
	// ContentTypeDefault is the default for text documents
	ContentTypeDefault = "text/plain"
	// ContentTypeJSON is the default for JSON strings
	ContentTypeJSON = "application/json"
)

func Body(body io.Reader) RequestMutation {
	return func(request *http.Request) error {
		if body == nil {
			request.Body = nil
			return nil
		}

		// If the ContentLength is actually 0, we can signal that ContentLength is ACTUALLY 0, and not just unknown
		OptimizeIfEmpty := func(request *http.Request) {
			if request.ContentLength == 0 {

				// This signals that the ContentLengt
				request.Body = http.NoBody
				request.GetBody = func() (io.ReadCloser, error) {
					return http.NoBody, nil
				}
			}
		}

		switch v := body.(type) {
		case *bytes.Buffer:
			request.ContentLength = int64(v.Len())
			buf := v.Bytes()
			request.GetBody = func() (io.ReadCloser, error) {
				r := bytes.NewReader(buf)
				return ioutil.NopCloser(r), nil
			}
			OptimizeIfEmpty(request)

		case *bytes.Reader:
			request.ContentLength = int64(v.Len())
			snapshot := *v
			request.GetBody = func() (io.ReadCloser, error) {
				r := snapshot
				return ioutil.NopCloser(&r), nil
			}
			OptimizeIfEmpty(request)

		case *strings.Reader:
			request.ContentLength = int64(v.Len())
			snapshot := *v
			request.GetBody = func() (io.ReadCloser, error) {
				r := snapshot
				return ioutil.NopCloser(&r), nil
			}
			OptimizeIfEmpty(request)

		default:
			// We don't have the same backwards compatibility issues as request.go:
			// https://github.com/golang/go/blob/226651a541/src/net/http/request.go#L823-L839

			// Convert the body to a ReadCloser
			rc, ok := body.(io.ReadCloser)
			if !ok {
				rc = ioutil.NopCloser(body)
			}
			request.Body = rc
			request.GetBody = func() (io.ReadCloser, error) {
				return rc, nil
			}

			// We don't know how to get the Length, this is signalled by setting ContentLength to 0
			request.ContentLength = 0
		}

		return nil
	}
}

// BodyFromJSON marshals the interface `v` into JSON for the request
func BodyFromJSON(v interface{}) RequestMutation {
	return func(request *http.Request) error {
		b := new(bytes.Buffer)
		if err := json.NewEncoder(b).Encode(v); err != nil {
			return errors.Wrap(err, "could not encode body")
		}
		if err := Body(b)(request); err != nil {
			return errors.Wrap(err, "could not set body")
		}
		request.Header.Set("Content-Type", ContentTypeJSON)
		return nil
	}
}

// BodyFromJSONString uses the string as JSON for the request
func BodyFromJSONString(s string) RequestMutation {
	return func(request *http.Request) error {
		if err := json.Unmarshal([]byte(s), &map[string]interface{}{}); err != nil {
			return errors.Wrap(err, "BodyFromJSONString received string that was not valid JSON")
		}
		if err := Body(strings.NewReader(s))(request); err != nil {
			return errors.Wrap(err, "could not set body")
		}
		request.Header.Set("Content-Type", ContentTypeJSON)
		return nil
	}
}

// WithHeader adds a header to the request
func WithHeader(key, value string) RequestMutation {
	return func(request *http.Request) error {
		request.Header.Add(key, value)
		return nil
	}
}

// BaseURL sets the URL of the request from a URL string
func BaseURL(base string) RequestMutation {
	return func(req *http.Request) error {
		u, err := url.Parse(base)
		if err != nil {
			return err
		}
		req.URL = u
		return nil
	}
}

// ResolvePath sets the path of the request by resolving it on the request.URL
func ResolvePath(path string) RequestMutation {
	return func(req *http.Request) error {
		u, err := url.Parse(path)
		if err != nil {
			return err
		}
		req.URL = req.URL.ResolveReference(u)
		return nil
	}
}
