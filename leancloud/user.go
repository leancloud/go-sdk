package leancloud

type User struct {
	Object
	SessionToken string `json:"sessionToken"`
}
