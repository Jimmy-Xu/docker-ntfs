package container

import (
	"golang.org/x/sys/unix"
	"github.com/Sirupsen/logrus"
)

func detachMounted(path string) error {
	logrus.Debugf("[detachMounted] Before - path:%v", path)
	return unix.Unmount(path, unix.MNT_DETACH)
}
