package leancloud

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"reflect"

	"github.com/levigross/grequests"
)

type CloudEngine interface {
	Init(client *Client)
	Handler() http.Handler

	Run(name string, params interface{}, runOptions ...RunOption) (interface{}, error)
	RPC(name string, params interface{}, results interface{}, runOptions ...RunOption) error

	Define(name string, fn func(*FunctionRequest) (interface{}, error), defineOptions ...DefineOption)

	BeforeSave(class string, fn func(*ClassHookRequest) (interface{}, error))
	AfterSave(class string, fn func(*ClassHookRequest) error)
	BeforeUpdate(class string, fn func(*ClassHookRequest) (interface{}, error))
	AfterUpdate(class string, fn func(*ClassHookRequest) error)
	BeforeDelete(class string, fn func(*ClassHookRequest) (interface{}, error))
	AfterDelete(class string, fn func(*ClassHookRequest) error)
	OnVerified(verifyType string, fn func(*ClassHookRequest) error)
	OnLogin(fn func(*ClassHookRequest) error)

	OnIMMessageReceived(fn func(*RealtimeHookRequest) (interface{}, error))
	OnIMReceiversOffline(fn func(*RealtimeHookRequest) (interface{}, error))
	OnIMMessageSent(fn func(*RealtimeHookRequest) error)
	OnIMMessageUpdate(fn func(*RealtimeHookRequest) (interface{}, error))
	OnIMConversationStart(fn func(*RealtimeHookRequest) (interface{}, error))
	OnIMConversationStarted(fn func(*RealtimeHookRequest) error)
	OnIMConversationAdd(fn func(*RealtimeHookRequest) (interface{}, error))
	OnIMConversationRemove(fn func(*RealtimeHookRequest) (interface{}, error))
	OnIMConversationAdded(fn func(*RealtimeHookRequest) error)
	OnIMConversationRemoved(fn func(*RealtimeHookRequest) error)
	OnIMConversationUpdate(fn func(*RealtimeHookRequest) (interface{}, error))
	OnIMClientOnline(fn func(*RealtimeHookRequest) error)
	OnIMClientOffline(fn func(*RealtimeHookRequest) error)

	client() *Client
}

type engine struct {
	c         *Client
	functions map[string]*functionType
}

var Engine CloudEngine

// FunctionRequest contains request information of Cloud Function
type FunctionRequest struct {
	Params       interface{}
	CurrentUser  *User
	SessionToken string
	Meta         map[string]string
}

// DefineOption apply options for definition of Cloud Function
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

// WithoutFetchUser don't fetch current user originated the request
func WithoutFetchUser() DefineOption {
	return &defineOption{
		fetchUser: false,
	}
}

// WithInternal restricts that the Cloud Function can only be executed in LeanEngine
func WithInternal() DefineOption {
	return &defineOption{
		internal: true,
	}
}

// RunOption apply options for execution of Cloud Function
type RunOption interface {
	apply(*map[string]interface{})
}

type runOption struct {
	remote       bool
	rpc          bool
	engine       *engine
	user         *User
	sessionToken string
}

func (option *runOption) apply(runOption *map[string]interface{}) {
	if option.remote {
		(*runOption)["remote"] = true
	}

	if option.user != nil {
		(*runOption)["user"] = option.user
	}

	if option.sessionToken != "" {
		(*runOption)["sessionToken"] = option.sessionToken
	}

	if option.engine != nil {
		(*runOption)["engine"] = option.engine
	}

	if option.rpc {
		(*runOption)["rpc"] = option.rpc
	}
}

// WithUser specifics the user of the calling
func WithUser(user *User) RunOption {
	return &runOption{
		user: user,
	}
}

// WithSessionToken specifics the sessionToken of the calling
func WithSessionToken(token string) RunOption {
	return &runOption{
		sessionToken: token,
	}
}

type functionType struct {
	call         func(*FunctionRequest) (interface{}, error)
	defineOption map[string]interface{}
}

func init() {
	Engine = &engine{
		functions: make(map[string]*functionType),
	}
}

// Init the LeanEngine part of Go SDK
func (engine *engine) Init(client *Client) {
	engine.c = client
}

func (engine *engine) client() *Client {
	if engine.c == nil {
		err := errors.New("not initialized (call leancloud.Engine.Init before use LeanEngine features)")
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		panic(err)
	}

	return engine.c
}

// Define declares a Cloud Function with name & options of definition
func (engine *engine) Define(name string, fn func(*FunctionRequest) (interface{}, error), defineOptions ...DefineOption) {
	if engine.functions[name] != nil {
		panic(fmt.Errorf("%s alreay defined", name))
	}

	engine.functions[name] = new(functionType)
	engine.functions[name].defineOption = map[string]interface{}{
		"fetchUser": true,
		"internal":  false,
	}

	for _, v := range defineOptions {
		v.apply(engine.functions[name])
	}

	engine.functions[name].call = fn
}

// Call cloud function locally
func (engine *engine) Run(name string, params interface{}, runOptions ...RunOption) (interface{}, error) {
	return callCloudFunction(engine.client(), name, params, (append(runOptions, &runOption{
		engine: engine,
	}))...)
}

// Call cloud function locally, bind response into `results`
func (engine *engine) RPC(name string, params interface{}, results interface{}, runOptions ...RunOption) error {
	response, err := callCloudFunction(engine.client(), name, params, (append(runOptions, &runOption{
		rpc:    true,
		engine: engine,
	}))...)

	if err := bind(reflect.Indirect(reflect.ValueOf(response)), reflect.Indirect(reflect.ValueOf(results))); err != nil {
		return nil
	}

	return err
}

// Call cloud funcion remotely
func (client *Client) Run(name string, params interface{}, runOptions ...RunOption) (interface{}, error) {
	return callCloudFunction(client, name, params, (append(runOptions, &runOption{
		remote: true,
	}))...)
}

// Call cloud function remotely, bind response into `results`
func (client *Client) RPC(name string, params interface{}, results interface{}, runOptions ...RunOption) error {
	response, err := callCloudFunction(client, name, params, (append(runOptions, &runOption{
		rpc:    true,
		remote: true,
	}))...)

	if err := bind(reflect.Indirect(reflect.ValueOf(response)), reflect.Indirect(reflect.ValueOf(results))); err != nil {
		return nil
	}

	return err
}

func callCloudFunction(client *Client, name string, params interface{}, runOptions ...RunOption) (interface{}, error) {
	options := make(map[string]interface{})
	sessionToken := ""
	var currentUser *User

	for _, v := range runOptions {
		v.apply(&options)
	}

	isRpc := options["rpc"] != nil && options["rpc"] != false

	if options["sessionToken"] != nil && options["user"] != nil {
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
		var path string

		reqOptions := client.getRequestOptions()

		if isRpc {
			path = fmt.Sprint("/1.1/call/", name)
			reqOptions.JSON = encode(params, true)
		} else {
			path = fmt.Sprint("/1.1/functions/", name)
			reqOptions.JSON = params
		}

		if sessionToken != "" {
			resp, err = client.request(methodPost, path, reqOptions, UseSessionToken(sessionToken))
		} else if currentUser != nil {
			resp, err = client.request(methodPost, path, reqOptions, UseUser(currentUser))
		} else {
			resp, err = client.request(methodPost, path, reqOptions)
		}

		if err != nil {
			return nil, err
		}

		respJSON := &functionResponse{}

		if err := json.Unmarshal(resp.Bytes(), respJSON); err != nil {
			return nil, err
		}

		if isRpc {
			return decode(respJSON.Result)
		} else {
			return respJSON.Result, err
		}
	}

	if options["engine"].(*engine).functions[name] == nil {
		return nil, fmt.Errorf("no such cloud function %s", name)
	}

	request := FunctionRequest{
		Params: params,
		Meta: map[string]string{
			"remoteAddr": "",
		},
	}

	if sessionToken != "" {
		request.SessionToken = sessionToken
		user, err := client.Users.Become(sessionToken)
		if err != nil {
			return nil, err
		}
		request.CurrentUser = user
	} else if currentUser != nil {
		request.CurrentUser = currentUser
		request.SessionToken = currentUser.SessionToken
	}

	return options["engine"].(*engine).functions[name].call(&request)
}
