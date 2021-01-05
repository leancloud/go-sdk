package leancloud

import (
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/levigross/grequests"
)

type function func(*Request) (interface{}, error)

type Request struct {
	Params       interface{}
	CurrentUser  *User
	SessionToken string
	Meta         map[string]string
}

type DefineOption interface {
	apply(*functionType)
}

type defineOption struct {
	fetchUser bool
	internal  bool
}

func (option *defineOption) apply(fn *functionType) {
	if option.fetchUser == false {
		fn.defineOption["fetchUser"] = false
	}

	if option.internal == true {
		fn.defineOption["internal"] = true
	}
}

func WithoutFetchUser() DefineOption {
	return &defineOption{
		fetchUser: false,
	}
}

func WithInteral() DefineOption {
	return &defineOption{
		internal: true,
	}
}

type RunOption interface {
	apply(*map[string]interface{})
}

type runOption struct {
	remote       bool
	user         *User
	sessionToken string
}

func (option *runOption) apply(runOption *map[string]interface{}) {
	if option.remote == true {
		(*runOption)["remote"] = true
	}

	if option.user != nil {
		(*runOption)["user"] = option.user
	}

	if option.sessionToken != "" {
		(*runOption)["sessionToken"] = option.sessionToken
	}
}

func WithRemote() RunOption {
	return &runOption{
		remote: true,
	}
}

func WithUser(user *User) RunOption {
	return &runOption{
		user: user,
	}
}

func WithSessionToken(token string) RunOption {
	return &runOption{
		sessionToken: token,
	}
}

type functionType struct {
	call         function
	defineOption map[string]interface{}
}

type functionError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

var client *Client

var functions map[string]*functionType

func init() {
	functions = make(map[string]*functionType)
	client = NewEnvClient()
}

func Define(name string, fn function, defineOptions ...DefineOption) {
	if functions[name] != nil {
		panic(fmt.Errorf("%s alreay defined", name))
	}

	functions[name] = new(functionType)
	functions[name].defineOption = map[string]interface{}{
		"fetchUser": true,
		"internal":  false,
	}

	for _, v := range defineOptions {
		v.apply(functions[name])
	}

	functions[name].call = fn
}

func Run(name string, object interface{}, runOptions ...RunOption) (interface{}, error) {
	options := make(map[string]interface{})
	sessionToken := ""
	var currentUser *User

	for _, v := range runOptions {
		v.apply(&options)
	}

	if options["sessionToken"] != "" && options["user"] != nil {
		return nil, fmt.Errorf("unable to set both of sessionToken & User")
	}

	if options["sessionToken"] != nil {
		sessionToken = options["sessionToken"].(string)
	}

	if options["user"] != nil {
		currentUser = options["user"].(*User)
	}

	if options["remote"] == true {
		var err error
		var resp *grequests.Response
		path := fmt.Sprint("/1.1/functions/", name)
		reqOption := client.getRequestOptions()
		reqOption.JSON = object
		if sessionToken != "" {
			resp, err = client.request(ServiceAPI, methodPost, path, reqOption, UseSessionToken(sessionToken))
		} else if currentUser != nil {
			resp, err = client.request(ServiceAPI, methodPost, path, reqOption, UseUser(currentUser))
		} else {
			resp, err = client.request(ServiceAPI, methodPost, path, reqOption)
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

	if functions[name] == nil {
		return nil, fmt.Errorf("no such cloud function %s", name)
	}

	request := Request{
		Params: object,
		Meta: map[string]string{
			"remoteAddr": "127.0.0.1",
		},
	}

	if sessionToken != "" {
		request.SessionToken = sessionToken
		user, err := client.Users.Become(sessionToken)
		if err != nil {
			return nil, err
		}
		request.CurrentUser = user
	}

	if currentUser != nil {
		request.CurrentUser = currentUser
		request.SessionToken = currentUser.SessionToken
	}

	return functions[name].call(&request)
}

func RPC(name string, object interface{}, ret interface{}, runOptions ...RunOption) (interface{}, error) {
	options := make(map[string]interface{})
	sessionToken := ""
	var currentUser *User

	for _, v := range runOptions {
		v.apply(&options)
	}

	if options["sessionToken"] != nil && options["user"] != nil {
		return nil, fmt.Errorf("unable to set both of sessionToken & User")
	}

	if options["sessionToken"] != nil {
		sessionToken = options["sessionToken"].(string)
	}

	if options["usr"] != nil {
		currentUser = options["user"].(*User)
	}

	if options["remote"] == true {
		var err error
		var resp *grequests.Response
		path := fmt.Sprint("/1.1/call/", name)
		reqOption := client.getRequestOptions()
		reqOption.JSON = encode(object, true)
		if sessionToken != "" {
			resp, err = client.request(ServiceAPI, methodPost, path, reqOption, UseSessionToken(sessionToken))
		} else if currentUser != nil {
			resp, err = client.request(ServiceAPI, methodPost, path, reqOption, UseUser(currentUser))
		} else {
			resp, err = client.request(ServiceAPI, methodPost, path, reqOption)
		}

		if err != nil {
			return nil, err
		}

		respJSON := new(functionResponse)
		if err := json.Unmarshal(resp.Bytes(), respJSON); err != nil {
			return nil, err
		}

		res, err := decode(respJSON.Result)
		if err != nil {
			return res, nil
		}

		if err := bind(reflect.Indirect(reflect.ValueOf(res)), reflect.Indirect(reflect.ValueOf(ret))); err != nil {
			return res, nil
		}

		return nil, nil
	}

	if functions[name] == nil {
		return nil, fmt.Errorf("no such cloud function %s", name)
	}

	request := Request{
		Params: object,
		Meta: map[string]string{
			"remoteAddr": "127.0.0.1",
		},
	}

	if sessionToken != "" {
		request.SessionToken = sessionToken
		user, err := client.Users.Become(sessionToken)
		if err != nil {
			return nil, err
		}
		request.CurrentUser = user
	}

	if currentUser != nil {
		request.CurrentUser = currentUser
		request.SessionToken = currentUser.SessionToken
	}

	res, err := functions[name].call(&request)
	if err != nil {
		return nil, err
	}

	if err := bind(reflect.Indirect(reflect.ValueOf(res)), reflect.Indirect(reflect.ValueOf(ret))); err != nil {
		return res, nil
	}

	return nil, nil
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
