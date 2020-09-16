package engine

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/leancloud/go-sdk/leancloud"
	"github.com/levigross/grequests"
)

type function func(*Request) (interface{}, error)

type Request struct {
	Params       interface{}
	CurrentUser  *leancloud.User
	SessionToken string
	Meta         http.Request
}

type DefineOption struct {
	NotFetchUser bool
	Internal     bool
}

type RunOption struct {
	Remote       bool
	User         *leancloud.User
	SessionToken string
}

type functionType struct {
	call function
	DefineOption
}

var client *leancloud.Client
var functions map[string]*functionType

func init() {
	functions = make(map[string]*functionType)
	client = leancloud.NewEnvClient()
}

func Define(name string, fn function) error {
	return define(name, fn, nil)
}

func DefineWithOption(name string, fn function, option *DefineOption) error {
	return define(name, fn, option)
}

func define(name string, fn function, option *DefineOption) error {
	if functions[name] != nil {
		return fmt.Errorf("%s alreay defined", name)
	}

	functions[name] = new(functionType)

	if option != nil {
		if !option.NotFetchUser {
			functions[name].NotFetchUser = false
		}
		if option.Internal {
			functions[name].Internal = true
		}
	} else {
		functions[name].DefineOption = DefineOption{
			NotFetchUser: true,
			Internal:     false,
		}
	}

	functions[name].call = fn

	return nil
}

func Run(name string, payload interface{}) (interface{}, error) {
	return run(name, payload, nil)
}

func RunWithOption(name string, payload interface{}, option *RunOption) (interface{}, error) {
	return run(name, payload, option)
}

func run(name string, payload interface{}, options *RunOption) (interface{}, error) {
	if options == nil {
		return runRemote(name, payload, options)
	}

	if options.SessionToken != "" && options.User != nil {
		return nil, fmt.Errorf("unable to set both of sessionToken & User")
	}

	if options.Remote {
		return runRemote(name, payload, options)
	}

	return runLocal(name, payload, options)
}

func runLocal(name string, payload interface{}, options *RunOption) (interface{}, error) {
	request := Request{
		Params: payload,
		Meta: http.Request{
			RemoteAddr: "127.0.0.1",
		},
	}
	if options.SessionToken != "" {
		user, err := client.Users.Become(options.SessionToken)
		if err != nil {
			return nil, err
		}
		request.CurrentUser = user
		request.SessionToken = options.SessionToken
	} else if options.User != nil {
		request.CurrentUser = options.User
		request.SessionToken = options.User.GetSessionToken()
	}
	return functions[name].call(&request)
}

func runRemote(name string, payload interface{}, options *RunOption) (interface{}, error) {
	var resp *grequests.Response
	var err error
	path := fmt.Sprint("/1.1/functions/", name)
	option := client.GetRequestOptions()
	option.JSON = payload
	if options == nil {
		resp, err = client.Request(leancloud.ServiceAPI, leancloud.MethodPost, path, option)
	} else {
		if options.SessionToken != "" {
			resp, err = client.Request(leancloud.ServiceAPI, leancloud.MethodPost, path, option, leancloud.UseSessionToken(options.SessionToken))
		} else if options.User != nil {
			resp, err = client.Request(leancloud.ServiceAPI, leancloud.MethodPost, path, option, leancloud.UseUser(options.User))
		}
	}
	if err != nil {
		return nil, err
	}
	respJSON := new(map[string]interface{})
	if err := json.Unmarshal(resp.Bytes(), respJSON); err != nil {
		return nil, err
	}
	return respJSON, err
}

func RPC(name string, payload interface{}) (interface{}, error) {
	return rpc(name, payload, nil)
}

func RPCWithOption(name string, payload interface{}, option *RunOption) (interface{}, error) {
	return rpc(name, payload, option)
}

func rpc(name string, payload interface{}, option *RunOption) (interface{}, error) {
	return nil, nil
}
