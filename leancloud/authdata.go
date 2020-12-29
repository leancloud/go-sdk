package leancloud

type AuthData struct {
	data map[string]map[string]interface{}
}

func NewAuthData() *AuthData {
	auth := new(AuthData)
	auth.data = make(map[string]map[string]interface{})
	return auth
}

func (auth *AuthData) Set(provider string, data map[string]interface{}) {
	auth.data[provider] = data
}

func (auth *AuthData) SetAnonymous(data map[string]interface{}) {
	auth.data["anonymous"] = data
}

func (auth *AuthData) Get(provider string) map[string]interface{} {
	return auth.data[provider]
}
