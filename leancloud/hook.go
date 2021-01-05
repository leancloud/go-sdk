package leancloud

import (
	"fmt"
	"os"
)

func hookAuthenticate(key string) bool {
	if key != os.Getenv("LEANCLOUD_APP_HOOK_KEY") {
		return false
	}

	return true
}

func defineHook(class, hook string, fn Function) {
	name := fmt.Sprint(hook, class)
	if functions[name] != nil {
		panic(fmt.Errorf("LeanEngine: %s of %s already defined", hook, class))
	}

	functions[name] = new(functionType)
	functions[name].defineOption = map[string]interface{}{
		"fetchUser": true,
		"internal":  false,
	}
	functions[name].call = fn
}

// BeforeSave will be called before saving an Object
func BeforeSave(class string, fn Function) {
	defineHook(class, "__before_save_for_", fn)
}

// AfterSave will be called after Object saved
func AfterSave(class string, fn Function) {
	defineHook(class, "__after_save_for_", fn)
}

// BeforeUpdate will be called before updating an Object
func BeforeUpdate(class string, fn Function) {
	defineHook(class, "__before_update_for_", fn)
}

// AfterUpdate will be called after Object updated
func AfterUpdate(class string, fn Function) {
	defineHook(class, "__after_update_for_", fn)
}

// BeforeDelete will be called before deleting an Object
func BeforeDelete(class string, fn Function) {
	defineHook(class, "__before_delete_for_", fn)
}

// AfterDelete will be called after Object deleted
func AfterDelete(class string, fn Function) {
	defineHook(class, "__after_delete_for_", fn)
}

// OnVerified will be called
func OnVerified(verifyType string, fn Function) {
	Define(fmt.Sprint("__on_verified_", verifyType), fn)
}

// OnLogin will be called
func OnLogin(fn Function) {
	Define("__on_login__User", fn)
}

func OnIMMessageReceived(fn Function) {

}

func OnIMReceiversOffline(fn Function) {

}

func OnIMMessageSent(fn Function) {

}

func OnIMMessageUpdate(fn Function) {

}

func OnImConversationStart(fn Function) {

}

func OnImConversationStarted(fn Function) {

}

func OnIMConversationAdd(fn Function) {

}

func OnIMConversationRemove(fn Function) {

}

func OnIMConversationAdded(fn Function) {

}

func OnIMConversationRemoved(fn Function) {

}

func OnIMConversationUpdate(fn Function) {

}

func OnIMClientOnline(fn Function) {

}

func OnIMClientOffline(fn Function) {

}
