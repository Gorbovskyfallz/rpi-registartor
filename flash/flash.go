package flash

import (
	flashInfo "github.com/moby/sys/mountinfo"
	"golang.org/x/sys/unix"
	"io"
	"log"
	"os/exec"
	"strings"
)

type StatFlash struct {
	Mounted     bool   // from checker
	MountPoint  string // from config
	DeviceName  string // from config
	FlashUse    FlashUse
	FlashErrors FlashErrors
}

type FlashErrors struct {
	MountExitStatus   int
	UnmountExitStatus int
}

type FlashUse struct {
	ServiceWork bool // check service selected from config
	ProcessWork bool // check process from config
}

// MountInfo "Checking - mounted or not to mountpath
// configMountPoint (/media/passed3/flash for e.x.)"
func (s *StatFlash) MountInfo(mountPoint string) (status bool, mountErr error) {
	name := "MountInfo"
	if status, mountErr = flashInfo.Mounted(mountPoint); mountErr != nil {
		log.Printf("%s: error \"%v\" occured\n", mountErr, name)
		return status, mountErr
	}
	s.Mounted = status
	return s.Mounted, mountErr
}

// MountFlash mounting block device dev to mountpoint with path
func (fe *FlashErrors) MountFlash(dev, mountPath string) {
	name := "MountFlash"
	var magicFlag uintptr = unix.MS_MGC_VAL
	err := unix.Mount(dev, mountPath, "exfat", magicFlag, "")
	if err != nil {
		switch {
		case err.Error() == "no such device":
			log.Printf("%s: device on mountPath %s\n", name, err)
			fe.MountExitStatus = 1
		case err.Error() == "no such file or directory":
			log.Printf("%s: mountpath: %s\n", name, err)
			fe.MountExitStatus = 2
		case err.Error() == "device or resource busy":
			log.Printf("%s: device: %s\n", name, err)
			fe.MountExitStatus = 3
		case err.Error() == "invalid argument":
			log.Printf("%s: arguments: %s\n", name, err)
			fe.MountExitStatus = 4
		default:
			log.Printf("%s: %s mnt to %s\n", name, dev, mountPath)
			fe.MountExitStatus = 0

		}

	}
	return
}

// UmountPoint unmount all flash from mediamountdir
func (fe *FlashErrors) UmountPoint(mountPoint string) int {
	name := "UmountPoint"
	if unmountErr := unix.Unmount(mountPoint, 0); unmountErr != nil {
		log.Printf("%s:error \"%s\" occured\n", name, unmountErr)
		fe.UnmountExitStatus = 1
	} else {
		log.Printf("%s: %s unmounted\n", name, mountPoint)
		fe.UnmountExitStatus = 0
	}

	return 0

}

// CheckPid check potentional disaster procces using the flash
// переделать на сисколы
func (f *FlashUse) CheckPid(processName string) bool {
	util := "pidof"
	args := "-s"
	out, _ := exec.Command(util, args, processName).Output()
	if len(out) != 0 {
		f.ProcessWork = true
	} else {
		f.ProcessWork = false
	}
	return f.ProcessWork
}

// CheckService checkin service that's processes ffmpeg+gpio things
func (f *FlashUse) CheckService(serviceName string) bool {
	// возможно, это стоит переписать на системных
	// вызовах без использованися exec
	// sudo systemctl status docker.service | grep Active
	name := "CheckSrvice"
	grep := exec.Command("grep", "Active")
	command := exec.Command("systemctl", "status", serviceName)
	pipe, _ := command.StdoutPipe()
	defer func(pipe io.ReadCloser) {
		closePipeErr := pipe.Close()
		if closePipeErr != nil {
			log.Printf("%s: %v", name, closePipeErr)
		}
	}(pipe)
	grep.Stdin = pipe
	startErr := command.Start()
	if startErr != nil {
		log.Printf("%s: %s command: %v", name, grep, startErr)
	}
	res, _ := grep.Output()
	if strings.Contains(string(res), "active (running)") {
		f.ServiceWork = true
	} else {
		f.ServiceWork = false
	}

	return f.ServiceWork
}
