package leancloud

import (
	"fmt"
	"os"
)

var storageHookmap = map[string]string{
	"beforeSave":   "__before_save_for_",
	"afterSave":    "__after_save_for_",
	"beforeUpdate": "__before_update_for_",
	"afterUpdate":  "__after_update_for_",
	"beforeDelete": "__before_delete_for_",
	"afterDelete":  "__after_delete_for_",
	"onVerified":   "__on_verified_",
	"onLogin":      "__on_login_",
}

var realtimeHookmap = map[string]string{
	"onIMMessageReceived":     "_messageReceived",
	"onIMReceiversOffline":    "_receiversOffline",
	"onIMMessageSent":         "_messageSent",
	"onIMMessageUpdate":       "_messageUpdate",
	"onIMConversationStart":   "_conversationStart",
	"onIMConversationStarted": "_conversationStarted",
	"onIMConversationAdd":     "_conversationAdd",
	"onIMConversationAdded":   "_conversationAdded",
	"onIMConversationRemove":  "_conversationRemove",
	"onIMConversationRemoved": "_conversationRemoved",
	"onIMConversationUpdate":  "_conversationUpdate",
	"onIMClientOnline":        "_clientOnline",
	"onIMClientOffline":       "_clientOffline",
	"onIMClientSign":          "_rtmClientSign",
}

func hookAuthenticate(key string) bool {
	if key != os.Getenv("LEANCLOUD_APP_HOOK_KEY") {
		return false
	}

	return true
}

func defineHook(class, hook string, fn func(*Object, *User) (interface{}, error)) {
	name := fmt.Sprint(hook, class)
	if functions[name] != nil {
		panic(fmt.Errorf("LeanEngine: %s of %s already defined", hook, class))
	}

	functions[name] = new(functionType)
	functions[name].defineOption = map[string]interface{}{
		"fetchUser": true,
		"internal":  false,
	}
	functions[name].call = func(r *Request) (interface{}, error) {
		if r.Params != nil {
			params, ok := r.Params.(map[string]interface{})
			if !ok {
				return nil, fmt.Errorf("invalid request body")
			}
			object, err := decodeObject(params["object"])
			if err != nil {
				return nil, err
			}
			if params["user"] != nil {
				user, err := decodeUser(params["user"])
				if err != nil {
					return nil, err
				}
				return fn(object, user)
			}
			return fn(object, nil)
		}

		return nil, nil
	}
}

// BeforeSave will be called before saving an Object
func BeforeSave(class string, fn func(*Object, *User) (interface{}, error)) {
	defineHook(class, "__before_save_for_", fn)
}

// AfterSave will be called after Object saved
func AfterSave(class string, fn func(*Object, *User) (interface{}, error)) {
	defineHook(class, "__after_save_for_", fn)
}

// BeforeUpdate will be called before updating an Object
func BeforeUpdate(class string, fn func(*Object, *User) (interface{}, error)) {
	defineHook(class, "__before_update_for_", fn)
}

// AfterUpdate will be called after Object updated
func AfterUpdate(class string, fn func(*Object, *User) (interface{}, error)) {
	defineHook(class, "__after_update_for_", fn)
}

// BeforeDelete will be called before deleting an Object
func BeforeDelete(class string, fn func(*Object, *User) (interface{}, error)) {
	defineHook(class, "__before_delete_for_", fn)
}

// AfterDelete will be called after Object deleted
func AfterDelete(class string, fn func(*Object, *User) (interface{}, error)) {
	defineHook(class, "__after_delete_for_", fn)
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

func OnIMMessageReceived(fn func(map[string]interface{}) (interface{}, error)) {
	Define(realtimeHookmap["onIMMessageReceived"], func(r *Request) (interface{}, error) {
		params, ok := r.Params.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("invalid request body")
		}

		return fn(params)
	})
}

func OnIMReceiversOffline(fn func(map[string]interface{}) (interface{}, error)) {
	Define(realtimeHookmap["onIMMessageReceiversOffline"], func(r *Request) (interface{}, error) {
		params, ok := r.Params.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("invalid request body")
		}

		return fn(params)
	})
}

func OnIMMessageSent(fn func(map[string]interface{}) (interface{}, error)) {
	Define(realtimeHookmap["onIMMessageSent"], func(r *Request) (interface{}, error) {
		params, ok := r.Params.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("invalid request body")
		}

		return fn(params)
	})
}

func OnIMMessageUpdate(fn func(map[string]interface{}) (interface{}, error)) {
	Define(realtimeHookmap["onIMMessageUpdated"], func(r *Request) (interface{}, error) {
		params, ok := r.Params.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("invalid request body")
		}

		return fn(params)
	})
}

func OnImConversationStart(fn func(map[string]interface{}) (interface{}, error)) {
	Define(realtimeHookmap["onIMConversationStart"], func(r *Request) (interface{}, error) {
		params, ok := r.Params.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("invalid request body")
		}

		return fn(params)
	})
}

func OnImConversationStarted(fn func(map[string]interface{}) (interface{}, error)) {
	Define(realtimeHookmap["onIMConversationStarted"], func(r *Request) (interface{}, error) {
		params, ok := r.Params.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("invalid request body")
		}

		return fn(params)
	})
}

func OnIMConversationAdd(fn func(map[string]interface{}) (interface{}, error)) {
	Define(realtimeHookmap["onIMConversationAdd"], func(r *Request) (interface{}, error) {
		params, ok := r.Params.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("invalid request body")
		}

		return fn(params)
	})
}

func OnIMConversationRemove(fn func(map[string]interface{}) (interface{}, error)) {
	Define(realtimeHookmap["onIMConversationRemove"], func(r *Request) (interface{}, error) {
		params, ok := r.Params.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("invalid request body")
		}

		return fn(params)
	})
}

func OnIMConversationAdded(fn func(map[string]interface{}) (interface{}, error)) {
	Define(realtimeHookmap["onIMConversationAdded"], func(r *Request) (interface{}, error) {
		params, ok := r.Params.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("invalid request body")
		}

		return fn(params)
	})
}

func OnIMConversationRemoved(fn func(map[string]interface{}) (interface{}, error)) {
	Define(realtimeHookmap["onIMConversationRemoved"], func(r *Request) (interface{}, error) {
		params, ok := r.Params.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("invalid request body")
		}

		return fn(params)
	})
}

func OnIMConversationUpdate(fn func(map[string]interface{}) (interface{}, error)) {
	Define(realtimeHookmap["onIMConversationUpdated"], func(r *Request) (interface{}, error) {
		params, ok := r.Params.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("invalid request body")
		}

		return fn(params)
	})
}

func OnIMClientOnline(fn func(map[string]interface{}) (interface{}, error)) {
	Define(realtimeHookmap["onIMClientOnline"], func(r *Request) (interface{}, error) {
		params, ok := r.Params.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("invalid request body")
		}

		return fn(params)
	})
}

func OnIMClientOffline(fn func(map[string]interface{}) (interface{}, error)) {
	Define(realtimeHookmap["onIMClientOffline"], func(r *Request) (interface{}, error) {
		params, ok := r.Params.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("invalid request body")
		}

		return fn(params)
	})
}
