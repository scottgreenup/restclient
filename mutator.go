package restclient

import (
	"net/http"
)

type RequestMutator struct {
	mutations []RequestMutation
}

func NewRequestMutator(commonMutations ...RequestMutation) *RequestMutator {
	return &RequestMutator{
		mutations: commonMutations,
	}
}

func (rm *RequestMutator) Mutate(req *http.Request, mutations ...RequestMutation) (*http.Request, error) {
	mutations = append(rm.mutations, mutations...)
	for _, mutation := range mutations {
		var err error
		req, err = mutation(req)
		if err != nil {
			return nil, err
		}
	}
	return req, nil
}

func (rm *RequestMutator) NewRequest(mutations ...RequestMutation) (*http.Request, error) {
	req, err := http.NewRequest("", "", nil)

	// This should NEVER happen, else something is very broken
	if err != nil {
		panic(err)
	}

	req, err = rm.Mutate(req, mutations...)

	if err != nil {
		return nil, err
	}

	return req, nil
}
