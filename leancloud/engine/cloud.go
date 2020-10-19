package engine

import (
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/leancloud/go-sdk/leancloud"
	"github.com/levigross/grequests"
)

type function func(*Request) (interface{}, error)

type Request struct {
	Params       interface{}
	CurrentUser  *leancloud.User
	SessionToken string
	Meta         map[string]string
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

type functionError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

var client *leancloud.Client

var functions map[string]*functionType

func init() {
	functions = make(map[string]*functionType)
	client = leancloud.NewEnvClient()
}

func Define(name string, fn function) {
	define(name, fn, nil)
}

func DefineWithOption(name string, fn function, option *DefineOption) {
	define(name, fn, option)
}

func define(name string, fn function, option *DefineOption) {
	if functions[name] != nil {
		panic(fmt.Errorf("%s alreay defined", name))
	}

	functions[name] = new(functionType)

	if option != nil {
		functions[name].DefineOption = *option
	} else {
		functions[name].DefineOption = DefineOption{
			NotFetchUser: true,
			Internal:     false,
		}
	}

	functions[name].call = fn
}

func Run(name string, payload interface{}) (interface{}, error) {
	return run(name, payload, nil)
}

func RunWithOption(name string, payload interface{}, option *RunOption) (interface{}, error) {
	return run(name, payload, option)
}

func run(name string, payload interface{}, options *RunOption) (interface{}, error) {
	if options == nil {
		return runRemote(name, payload, nil)
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
		Meta: map[string]string{
			"remoteAddr": "127.0.0.1",
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
	if payload != nil {
		option.JSON = payload
	}
	if options == nil {
		resp, err = client.Request(leancloud.ServiceAPI, leancloud.MethodPost, path, option)
	} else {
		if options.SessionToken != "" {
			resp, err = client.Request(leancloud.ServiceAPI, leancloud.MethodPost, path, option, leancloud.UseSessionToken(options.SessionToken))
		} else if options.User != nil {
			resp, err = client.Request(leancloud.ServiceAPI, leancloud.MethodPost, path, option, leancloud.UseUser(options.User))
		} else {
			resp, err = client.Request(leancloud.ServiceAPI, leancloud.MethodPost, path, option)
		}
	}
	if err != nil {
		return nil, err
	}

	respJSON := new(functionResponse)
	if err := json.Unmarshal(resp.Bytes(), respJSON); err != nil {
		return nil, err
	}

	return respJSON.Result, err
}

func RPC(name string, payload interface{}) (interface{}, error) {
	return rpc(name, payload, nil)
}

func RPCWithOption(name string, payload interface{}, options *RunOption) (interface{}, error) {
	return rpc(name, payload, options)
}

func rpc(name string, payload interface{}, options *RunOption) (interface{}, error) {
	encodedPayload, err := encode(payload)
	if err != nil {
		return nil, err
	}

	resp, err := run(name, encodedPayload, options)
	if err != nil {
		return nil, err
	}

	object := new(leancloud.Object)
	decode(resp, object)
	if err := decode(resp, object); err != nil {
		return nil, err
	}

	return object, nil
}

func (ferr *functionError) Error() string {
	errString, err := json.Marshal(ferr)
	if err != nil {
		return fmt.Sprint(err)
	}

	return string(errString)
}

func ErrorWithCode(code int, message string) *functionError {
	return &functionError{
		Code:    code,
		Message: message,
	}
}

func Error(message string) *functionError {
	return &functionError{
		Code:    1,
		Message: message,
	}
}

func encode(payload interface{}) (interface{}, error) {
	payloadMap := new(map[string]interface{})
	payloadValue := reflect.ValueOf(payload)
	payloadType := payloadValue.Type()

	switch payloadType.Kind() {
	case reflect.Array:
		fallthrough
	case reflect.Slice:
	case reflect.Struct:
		return leancloud.EncodeObject(payload), nil
	default:
		return payload, nil
	}

	return payloadMap, nil
}

func decode(payload interface{}, object interface{}) error {
	payloadValue := reflect.ValueOf(payload)
	payloadType := payloadValue.Type()

	switch payloadType.Kind() {
	case reflect.Array:
		fallthrough
	case reflect.Slice:
	case reflect.Map:
		payloadMap, ok := payload.(map[string]interface{})
		if !ok {
			return fmt.Errorf("unexpected payload format for decoding")
		}
		return leancloud.DecodeObject(payloadMap, object)
	default:
		object = payload
		return nil
	}
	return nil
}
