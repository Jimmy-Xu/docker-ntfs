package mount

import (
	"syscall"
	"github.com/Sirupsen/logrus"
	"os/exec"
)

func mount(device, target, mType string, flag uintptr, data string) error {
	logrus.Infof("[mount] Begin device:%v, target:%v, mType:%v, flag:%v, data:%v", device, target, mType, flag, data)

	logrus.Infof("[mount] Before mount device:%v, target:%v, mType:%v, flag:%v, data:%v", device, target, mType, flag, data)
	if mType == "ntfs-3g" {
		args := []string{}
		args = append(args, device)
		args = append(args, target)
		if err := exec.Command("/sbin/mount.ntfs-3g", args...).Run(); err != nil {
			logrus.Infof("[mount] After mount.ntfs-3g err:%v", err)
			return err
		}
	} else if err := syscall.Mount(device, target, mType, flag, data); err != nil {
		logrus.Infof("[mount] After syscall.Mount() err:%v", err)
		return err
	}
	logrus.Infof("[mount] After mount device:%v, target:%v, mType:%v, flag:%v, data:%v", device, target, mType, flag, data)

	logrus.Infof("[mount] Before syscall.Mount() remount: device:%v, target:%v, mType:%v, flag:%v, data:%v", device, target, mType, flag, data)
	// If we have a bind mount or remount, remount...
	if flag&syscall.MS_BIND == syscall.MS_BIND && flag&syscall.MS_RDONLY == syscall.MS_RDONLY {
		//return syscall.Mount(device, target, mType, flag|syscall.MS_REMOUNT, data)
		if err := syscall.Mount(device, target, mType, flag|syscall.MS_REMOUNT, data); err != nil {
			logrus.Infof("[mount] After syscall.Mount() remount: err:%v", err)
			return err
		}
	}
	logrus.Infof("[mount] After syscall.Mount() remount: device:%v, target:%v, mType:%v, flag:%v, data:%v", device, target, mType, flag, data)

	logrus.Infof("[mount] End device:%v, target:%v, mType:%v, flag:%v, data:%v", device, target, mType, flag, data)
	return nil
}

func unmount(target string, flag int) error {
	return syscall.Unmount(target, flag)
}
