package engine

import (
	"net/http"

	"github.com/leancloud/go-sdk/leancloud"
)

type function func(*Request) error

type Request struct {
	Params       interface{}     `json:"params"`
	CurrentUser  *leancloud.User `json:"currentUser"`
	SessionToken *string         `json:"sessionToken"`
	Meta         interface{}     `json:"meta"`
}

type DefineOption struct {
	FetchUser bool `json:"fetchUser"`
	Internal  bool `json:"internal"`
}

type RunOption struct {
	Remote       *bool           `json:"remote"`
	User         *leancloud.User `json:"user"`
	SessionToken *string         `json:"sessionToken"`
	Req          *http.Request   `json:"req"`
}

func define(name string, fn function, options ...DefineOption) {

}

func run(name string, payload interface{}, options ...RunOption) (interface{}, error) {
	return nil, nil
}
