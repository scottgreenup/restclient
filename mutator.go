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

func (rm *RequestMutator) Mutate(req *http.Request, mutations ...RequestMutation) error {
	mutations = append(rm.mutations, mutations...)
	for _, mutation := range mutations {
		if err := mutation(req); err != nil {
			return err
		}
	}
	return nil
}

func (rm *RequestMutator) NewRequest(mutations ...RequestMutation) (*http.Request, error) {
	req, err := http.NewRequest("", "", nil)

	// This should NEVER happen, else something is very broken
	if err != nil {
		panic(err)
	}

	if err := rm.Mutate(req, mutations...); err != nil {
		return nil, err
	}

	return req, nil
}
