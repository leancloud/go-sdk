package leancloud

import (
	"fmt"
)

// ClassHookRequest contains object and user passed by Class hook calling
type ClassHookRequest struct {
	Object *Object
	User   *User
	Meta   map[string]string
}

// UpdatedKeys return keys which would be updated, only valid in beforeUpdate hook
func (r *ClassHookRequest) UpdatedKeys() []string {
	return r.Object.fields["_updatedKeys"].([]string)
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

func (engine *engine) defineClassHook(class, hook string, fn func(*ClassHookRequest) (interface{}, error)) {
	name := fmt.Sprint(hook, class)
	if engine.functions[name] != nil {
		panic(fmt.Errorf("LeanEngine: %s of %s already defined", hook, class))
	}

	engine.functions[name] = new(functionType)
	engine.functions[name].defineOption = map[string]interface{}{
		"fetchUser": true,
		"internal":  false,
		"hook":      true,
	}
	engine.functions[name].call = func(r *FunctionRequest) (interface{}, error) {
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
func (engine *engine) BeforeSave(class string, fn func(*ClassHookRequest) (interface{}, error)) {
	engine.defineClassHook(class, "__before_save_for_", fn)
}

// AfterSave will be called after Object saved
func (engine *engine) AfterSave(class string, fn func(*ClassHookRequest) error) {
	engine.defineClassHook(class, "__after_save_for_", func(r *ClassHookRequest) (interface{}, error) {
		return nil, fn(r)
	})
}

// BeforeUpdate will be called before updating an Object
func (engine *engine) BeforeUpdate(class string, fn func(*ClassHookRequest) (interface{}, error)) {
	engine.defineClassHook(class, "__before_update_for_", fn)
}

// AfterUpdate will be called after Object updated
func (engine *engine) AfterUpdate(class string, fn func(*ClassHookRequest) error) {
	engine.defineClassHook(class, "__after_update_for_", func(r *ClassHookRequest) (interface{}, error) {
		return nil, fn(r)
	})
}

// BeforeDelete will be called before deleting an Object
func (engine *engine) BeforeDelete(class string, fn func(*ClassHookRequest) (interface{}, error)) {
	engine.defineClassHook(class, "__before_delete_for_", fn)
}

// AfterDelete will be called after Object deleted
func (engine *engine) AfterDelete(class string, fn func(*ClassHookRequest) error) {
	engine.defineClassHook(class, "__after_delete_for_", func(r *ClassHookRequest) (interface{}, error) {
		return nil, fn(r)
	})
}

// OnVerified will be called when user was online
func (engine *engine) OnVerified(verifyType string, fn func(*ClassHookRequest) error) {
	engine.Define(fmt.Sprint("__on_verified_", verifyType), func(r *FunctionRequest) (interface{}, error) {
		params, ok := r.Params.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("invalid request body")
		}
		user, err := decodeUser(params["object"])
		if err != nil {
			return nil, err
		}
		req := ClassHookRequest{
			User: user,
			Meta: r.Meta,
		}
		return nil, fn(&req)
	})
}

// OnLogin will be called when user logged in
func (engine *engine) OnLogin(fn func(*ClassHookRequest) error) {
	engine.Define("__on_login__User", func(r *FunctionRequest) (interface{}, error) {
		params, ok := r.Params.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("invalid request body")
		}
		user, err := decodeUser(params["object"])
		if err != nil {
			return nil, err
		}
		req := ClassHookRequest{
			User: user,
			Meta: r.Meta,
		}
		return nil, fn(&req)
	})
}

func (engine *engine) defineRealtimeHook(name string, fn func(*RealtimeHookRequest) (interface{}, error)) {
	engine.Define(name, func(r *FunctionRequest) (interface{}, error) {
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
	engine.functions[name].defineOption["hook"] = true
}

func (engine *engine) OnIMMessageReceived(fn func(*RealtimeHookRequest) (interface{}, error)) {
	engine.defineRealtimeHook("_messageReceived", fn)
}

func (engine *engine) OnIMReceiversOffline(fn func(*RealtimeHookRequest) (interface{}, error)) {
	engine.defineRealtimeHook("_receiverOffline", fn)
}

func (engine *engine) OnIMMessageSent(fn func(*RealtimeHookRequest) error) {
	engine.defineRealtimeHook("_messageSent", func(r *RealtimeHookRequest) (interface{}, error) {
		return nil, fn(r)
	})
}

func (engine *engine) OnIMMessageUpdate(fn func(*RealtimeHookRequest) (interface{}, error)) {
	engine.defineRealtimeHook("_messageUpdate", fn)
}

func (engine *engine) OnIMConversationStart(fn func(*RealtimeHookRequest) (interface{}, error)) {
	engine.defineRealtimeHook("_conversationStart", fn)
}

func (engine *engine) OnIMConversationStarted(fn func(*RealtimeHookRequest) error) {
	engine.defineRealtimeHook("_conversationStarted", func(r *RealtimeHookRequest) (interface{}, error) {
		return nil, fn(r)
	})
}

func (engine *engine) OnIMConversationAdd(fn func(*RealtimeHookRequest) (interface{}, error)) {
	engine.defineRealtimeHook("_conversationStarted", fn)
}

func (engine *engine) OnIMConversationRemove(fn func(*RealtimeHookRequest) (interface{}, error)) {
	engine.defineRealtimeHook("_conversationRemove", fn)
}

func (engine *engine) OnIMConversationAdded(fn func(*RealtimeHookRequest) error) {
	engine.defineRealtimeHook("_conversationAdded", func(r *RealtimeHookRequest) (interface{}, error) {
		return nil, fn(r)
	})
}

func (engine *engine) OnIMConversationRemoved(fn func(*RealtimeHookRequest) error) {
	engine.defineRealtimeHook("_conversationRemoved", func(r *RealtimeHookRequest) (interface{}, error) {
		return nil, fn(r)
	})
}

func (engine *engine) OnIMConversationUpdate(fn func(*RealtimeHookRequest) (interface{}, error)) {
	engine.defineRealtimeHook("_conversationUpdate", fn)
}

func (engine *engine) OnIMClientOnline(fn func(*RealtimeHookRequest) error) {
	engine.defineRealtimeHook("_clientOnline", func(r *RealtimeHookRequest) (interface{}, error) {
		return nil, fn(r)
	})
}

func (engine *engine) OnIMClientOffline(fn func(*RealtimeHookRequest) error) {
	engine.defineRealtimeHook("_clientOffline", func(r *RealtimeHookRequest) (interface{}, error) {
		return nil, fn(r)
	})
}
