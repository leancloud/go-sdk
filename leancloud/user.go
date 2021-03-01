package leancloud

// User is a local representation of a user persisted to the LeanCloud server.
type User struct {
	Object
	SessionToken        string `json:"sessionToken"`
	Username            string `json:"username"`
	Email               string `json:"email"`
	EmailVerified       bool   `json:"emailVerified"`
	MobilePhoneNumber   string `json:"mobilePhoneNumber"`
	MobilePhoneVerified bool   `json:"mobilePhoneVerified"`
}
