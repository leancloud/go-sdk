package leancloud

import (
	"fmt"
	"os"
)

// ClassHookRequest contains object and user passed by Class hook calling
type ClassHookRequest struct {
	Object *Object
	User   *User
	Meta   map[string]string
}

// RealtimeHookRequest contains parameters passed by RTM hook calling
type RealtimeHookRequest struct {
	Params map[string]interface{}
	Meta   map[string]string
}

var classHookmap = map[string]string{
	"beforeSave":   "__before_save_for_",
	"afterSave":    "__after_save_for_",
	"beforeUpdate": "__before_update_for_",
	"afterUpdate":  "__after_update_for_",
	"beforeDelete": "__before_delete_for_",
	"afterDelete":  "__after_delete_for_",
	"onVerified":   "__on_verified_",
	"onLogin":      "__on_login_",
}

func hookAuthenticate(key string) bool {
	if key != os.Getenv("LEANCLOUD_APP_HOOK_KEY") {
		return false
	}

	return true
}

func defineClassHook(class, hook string, fn func(*ClassHookRequest) (interface{}, error)) {
	name := fmt.Sprint(hook, class)
	if functions[name] != nil {
		panic(fmt.Errorf("LeanEngine: %s of %s already defined", hook, class))
	}

	functions[name] = new(functionType)
	functions[name].defineOption = map[string]interface{}{
		"fetchUser": true,
		"internal":  false,
		"hook":      true,
	}
	functions[name].call = func(r *Request) (interface{}, error) {
		if r.Params != nil {
			req := new(ClassHookRequest)
			params, ok := r.Params.(map[string]interface{})
			if !ok {
				return nil, fmt.Errorf("invalid request body")
			}
			object, err := decodeObject(params["object"])
			if err != nil {
				return nil, err
			}
			req.Object = object
			if params["user"] != nil {
				user, err := decodeUser(params["user"])
				if err != nil {
					return nil, err
				}
				req.User = user
			}
			return fn(req)
		}

		return nil, nil
	}
}

// BeforeSave will be called before saving an Object
func BeforeSave(class string, fn func(*ClassHookRequest) (interface{}, error)) {
	defineClassHook(class, "__before_save_for_", fn)
}

// AfterSave will be called after Object saved
func AfterSave(class string, fn func(*ClassHookRequest) error) {
	defineClassHook(class, "__after_save_for_", func(r *ClassHookRequest) (interface{}, error) {
		return nil, fn(r)
	})
}

// BeforeUpdate will be called before updating an Object
func BeforeUpdate(class string, fn func(*ClassHookRequest) (interface{}, error)) {
	defineClassHook(class, "__before_update_for_", fn)
}

// AfterUpdate will be called after Object updated
func AfterUpdate(class string, fn func(*ClassHookRequest) error) {
	defineClassHook(class, "__after_update_for_", func(r *ClassHookRequest) (interface{}, error) {
		return nil, fn(r)
	})
}

// BeforeDelete will be called before deleting an Object
func BeforeDelete(class string, fn func(*ClassHookRequest) (interface{}, error)) {
	defineClassHook(class, "__before_delete_for_", fn)
}

// AfterDelete will be called after Object deleted
func AfterDelete(class string, fn func(*ClassHookRequest) error) {
	defineClassHook(class, "__after_delete_for_", func(r *ClassHookRequest) (interface{}, error) {
		return nil, fn(r)
	})
}

// OnVerified will be called when user was online
func OnVerified(verifyType string, fn func(*User) error) {
	Define(fmt.Sprint("__on_verified_", verifyType), func(r *Request) (interface{}, error) {
		params, ok := r.Params.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("invalid request body")
		}
		user, err := decodeUser(params["object"])
		if err != nil {
			return nil, err
		}

		return nil, fn(user)
	})
}

// OnLogin will be called when user logged in
func OnLogin(fn func(*User) error) {
	Define("__on_login__User", func(r *Request) (interface{}, error) {
		params, ok := r.Params.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("invalid request body")
		}
		user, err := decodeUser(params["object"])
		if err != nil {
			return nil, err
		}

		return nil, fn(user)
	})
}

func defineRealtimeHook(name string, fn func(*RealtimeHookRequest) (interface{}, error)) {
	Define(name, func(r *Request) (interface{}, error) {
		params, ok := r.Params.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("invalid request body")
		}
		req := RealtimeHookRequest{
			Params: params,
			Meta:   r.Meta,
		}
		return fn(&req)
	})
	functions[name].defineOption["hook"] = true
}

func OnIMMessageReceived(fn func(*RealtimeHookRequest) (interface{}, error)) {
	defineRealtimeHook("_messageReceived", fn)
}

func OnIMReceiversOffline(fn func(*RealtimeHookRequest) (interface{}, error)) {
	defineRealtimeHook("_receiverOffline", fn)
}

func OnIMMessageSent(fn func(*RealtimeHookRequest) error) {
	defineRealtimeHook("_messageSent", func(r *RealtimeHookRequest) (interface{}, error) {
		return nil, fn(r)
	})
}

func OnIMMessageUpdate(fn func(*RealtimeHookRequest) (interface{}, error)) {
	defineRealtimeHook("_messageUpdate", fn)
}

func OnImConversationStart(fn func(*RealtimeHookRequest) (interface{}, error)) {
	defineRealtimeHook("_conversationStart", fn)
}

func OnImConversationStarted(fn func(*RealtimeHookRequest) error) {
	defineRealtimeHook("_conversationStarted", func(r *RealtimeHookRequest) (interface{}, error) {
		return nil, fn(r)
	})
}

func OnIMConversationAdd(fn func(*RealtimeHookRequest) (interface{}, error)) {
	defineRealtimeHook("_conversationStarted", fn)
}

func OnIMConversationRemove(fn func(*RealtimeHookRequest) (interface{}, error)) {
	defineRealtimeHook("_conversationRemove", fn)
}

func OnIMConversationAdded(fn func(*RealtimeHookRequest) error) {
	defineRealtimeHook("_conversationAdded", func(r *RealtimeHookRequest) (interface{}, error) {
		return nil, fn(r)
	})
}

func OnIMConversationRemoved(fn func(*RealtimeHookRequest) error) {
	defineRealtimeHook("_conversationRemoved", func(r *RealtimeHookRequest) (interface{}, error) {
		return nil, fn(r)
	})
}

func OnIMConversationUpdate(fn func(*RealtimeHookRequest) (interface{}, error)) {
	defineRealtimeHook("_conversationUpdate", fn)
}

func OnIMClientOnline(fn func(*RealtimeHookRequest) error) {
	defineRealtimeHook("_clientOnline", func(r *RealtimeHookRequest) (interface{}, error) {
		return nil, fn(r)
	})
}

func OnIMClientOffline(fn func(*RealtimeHookRequest) error) {
	defineRealtimeHook("_clientOffline", func(r *RealtimeHookRequest) (interface{}, error) {
		return nil, fn(r)
	})
}
