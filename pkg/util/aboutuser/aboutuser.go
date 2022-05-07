// Package aboutuser
// 关于用户以及用户组的存在性检查、新建
package aboutuser

import (
	"fmt"
	"os/user"
)

//UserInfo 反馈给定用户在当前服务器上的用户信息
//反馈user.User 的结构体，包括
//type User struct {
//	Uid      string
//	Gid      string
//	Username string
//	Name     string
//	HomeDir  string
//}
func UserInfo(uname user.User) (*user.User, error) {
	lookup, err := user.Lookup(uname.Name)
	if err != nil {
		return &user.User{}, nil
	}
	fmt.Printf("%+v", *lookup)
	return lookup, nil
}

// IsUserExist 判断当前用户的uname在当前服务器上是否存在
// uname 用户的username
func IsUserExist(uname string) (bool, error) {
	lookup, err := user.Lookup(uname)
	if err != nil {
		return false, nil
	}
	fmt.Printf("%+v", *lookup)
	return true, nil
}

//IsUserGroupExist 判断用户组在当前机器上是否存在
func IsUserGroupExist(uname user.User) (bool, error) {
	lookup, err := user.LookupGroupId(uname.Gid)
	if err != nil {
		return false, nil
	}
	fmt.Printf("%+v", *lookup)
	return true, nil
}

// GetUserHomeDir 获取当前用户的家目录
func GetUserHomeDir(uname user.User) (string, error) {
	userExist, err := UserInfo(uname)
	if err != nil {
		return "", nil
	}
	return userExist.HomeDir, nil
}
