package leancloud

import "fmt"

// ACL include permission group of object
type ACL struct {
	content map[string]map[string]bool
}

func NewACL() *ACL {
	acl := new(ACL)
	acl.content = make(map[string]map[string]bool)
	return acl
}

func NewACLWithUser(user *User) *ACL {
	acl := NewACL()
	acl.set(user.ID, "read", true)
	acl.set(user.ID, "write", true)
	return acl
}

func (acl *ACL) SetPublicReadAccess(allowed bool) {
	acl.set("*", "read", allowed)
}

func (acl *ACL) SetPublicWriteAccess(allowed bool) {
	acl.set("*", "write", allowed)
}

func (acl *ACL) SetWriteAccess(user *User, allowed bool) {
	acl.set(user.ID, "write", allowed)
}

func (acl *ACL) SetReadAccess(user *User, allowed bool) {
	acl.set(user.ID, "read", allowed)
}

func (acl *ACL) SetRoleReadAccess(role *Role, allowed bool) {
	acl.set(fmt.Sprint("role:", role.Name), "read", allowed)
}

func (acl *ACL) SetRoleWriteAccess(role *Role, allowed bool) {
	acl.set(fmt.Sprint("role:", role.Name), "write", allowed)
}

func (acl *ACL) GetPublicReadAccess() bool {
	return acl.get("*", "read")
}

func (acl *ACL) GetPublicWriteAccess() bool {
	return acl.get("*", "write")
}

func (acl *ACL) GetReadAccess(user *User) bool {
	return acl.get(user.ID, "read")
}

func (acl *ACL) GetWriteAccess(user *User) bool {
	return acl.get(user.ID, "write")
}

func (acl *ACL) GetRoleReadAccess(role *Role) bool {
	return acl.get(fmt.Sprint("role:", role.Name), "read")
}

func (acl *ACL) GetRoleWriteAccess(role *Role) bool {
	return acl.get(fmt.Sprint("role:", role.Name), "write")
}

func (acl *ACL) set(key, perm string, allowed bool) {
	acl.content[key][perm] = allowed
}

func (acl *ACL) get(key, perm string) bool {
	return acl.content[key][perm]
}
