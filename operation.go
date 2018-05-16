package restclient

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/pkg/errors"
	"regexp"
)

const (
	// ContentTypeDefault is the default for text documents
	ContentTypeDefault = "text/plain"
	// ContentTypeJSON is the default for JSON strings
	ContentTypeJSON = "application/json"
)

// Operation represents a single request to a RESTful endpoint, it is used to build the actual http.Request
type Operation struct {
	config          *Config
	pathTemplate    string
	pathTemplateVar map[string]string
	body            io.Reader
	method          string
	contentType     string
	headers         map[string]string
	authentication  []interface{}
	merr            *multierror.Error
}

func newOperation(method string, config *Config) *Operation {
	return &Operation{
		pathTemplateVar: make(map[string]string),
		headers:         make(map[string]string),
		method:          method,
		contentType:     ContentTypeDefault,
		config:          config,
	}
}

// BodyFromJSON marshals the interface `v` into JSON for the request
func (o *Operation) BodyFromJSON(v interface{}) *Operation {
	b := new(bytes.Buffer)
	err := json.NewEncoder(b).Encode(v)

	if err != nil {
		o.merr = multierror.Append(o.merr, errors.Wrap(err, "could not encode body"))
		return o
	}

	o.body = b
	o.contentType = ContentTypeJSON
	return o
}

// BodyFromJSONString uses the string as JSON for the request
func (o *Operation) BodyFromJSONString(s string) *Operation {
	if err := json.Unmarshal([]byte(s), &map[string]interface{}{}); err != nil {
		o.merr = multierror.Append(o.merr, errors.Wrap(err, "BodyFromJSONString received string that was not valid JSON"))
	} else {
		o.body = strings.NewReader(s)
		o.contentType = ContentTypeJSON
	}
	return o
}

// WithPath changes the endpoint of the RESTful API that we are hitting. The path is treated as a template, variables
// are signalled with {}, e.g. "/boards/{id}"
func (o *Operation) WithPath(template string) *Operation {
	if o.pathTemplate != "" {
		o.merr = multierror.Append(o.merr, errors.New("WithPath was already called"))
		return o
	}

	if template == "" {
		o.merr = multierror.Append(o.merr, errors.New("WithPath was called with empty string"))
		return o
	}

	o.pathTemplate = template
	return o
}

// WithPathVar sets a variable for the path template
func (o *Operation) WithPathVar(key, value string) *Operation {
	o.pathTemplateVar[key] = value
	return o
}

// WithHeader adds a header to the request
func (o *Operation) WithHeader(key, value string) *Operation {
	o.headers[key] = value
	return o
}

// Authenticate activates the authentication method given to the config
func (o *Operation) Authenticate(v ...interface{}) *Operation {
	if o.config.AuthMethod == nil {
		o.merr = multierror.Append(o.merr, errors.New("Authenticate called without SetAuthenticationMethod"))
		return o
	}
	o.authentication = v
	return o
}

// BuildRequest builds the http.Request from the operation
func (o *Operation) BuildRequest() (*http.Request, error) {
	if o.merr != nil {
		return nil, o.merr
	}

	u, err := o.renderURL()
	if err != nil {
		return nil, errors.Wrap(err, "could not render URL")
	}

	request, err := http.NewRequest(o.method, u.String(), o.body)
	if err != nil {
		return nil, errors.Wrap(err, "could not create request")
	}

	for key, value := range o.headers {
		request.Header.Add(key, value)
	}

	request.Header.Add("Content-Type", o.contentType)

	if o.config.AuthMethod != nil {
		request, err = o.config.AuthMethod(request, o.authentication...)
		if err != nil {
			return nil, errors.Wrap(err, "could embed authentication into request")
		}
	}

	return request, nil
}

func (o *Operation) renderURL() (*url.URL, error) {
	o.pathTemplate = strings.TrimPrefix(o.pathTemplate, "/")

	path, err := format(o.pathTemplate, o.pathTemplateVar)
	if err != nil {
		return nil, err
	}

	u, err := url.Parse(path)
	if err != nil {
		return nil, err
	}

	return o.config.Endpoint.ResolveReference(u), nil
}

var (
	formatRegExp = regexp.MustCompile("{[A-Za-z0-9_]+}")
)

func format(format string, vars map[string]string) (string, error) {
	if vars == nil {
		vars = make(map[string]string)
	}

	used := map[string]bool{}
	for k := range vars {
		used[k] = false
	}

	b := []byte(format)
	loc := formatRegExp.FindIndex(b)

	newString := ""
	prevEnd := 0

	for loc != nil {
		key := string(b[loc[0]+1 : loc[1]-1])
		value, ok := vars[key]

		if !ok {
			return "", errors.Errorf("url contained reference to variable %s which was not supplied through WithPathVar", key)
		}

		if prevEnd == loc[0] {
			newString = newString + value
		} else {
			newString = newString + string(b[prevEnd:loc[0]]) + value
		}

		used[key] = true
		prevEnd = loc[1]

		if prevEnd < len(b) {
			loc = formatRegExp.FindIndex(b[prevEnd:])
			loc[0] += prevEnd
			loc[1] += prevEnd
		} else {
			loc = nil
		}
	}

	if prevEnd < len(b) {
		newString = newString + string(b[prevEnd:])
	}

	for key, wasUsed := range used {
		if !wasUsed {
			return "", errors.Errorf("url did not use %s which was provided through WithPathVar", key)
		}
	}

	return newString, nil

}
