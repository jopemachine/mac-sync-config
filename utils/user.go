package utils

import (
	"os"
	"os/user"
)

func IsRootUser() bool {
	return os.Geteuid() == 0
}

func GetMacosUserName() string {
	user, err := user.Current()
	PanicIfErr(err)
	return user.Username
}
