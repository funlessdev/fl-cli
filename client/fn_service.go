package client

import "net/http"

//					 default
// https://${HOST}/{NAMESPACE}/fn/hello
// https://${HOST}/{NAMESPACE}/ev/event1
// https://${HOST}/{NAMESPACE}/pkg/{PACKAGE}/fn/hello

type FnService struct {
	*Client
}

type Fn struct {
	Name string `json:"name,omitempty"`
}

func (fn *FnService) Invoke(fnName string) (*http.Response, error) {
	// https://${HOST}/{NAMESPACE}/fn/hello
	req, err := fn.CreateGet("fn/" + fnName)
	if err != nil {
		return nil, err
	}
	return fn.Send(req)
}
